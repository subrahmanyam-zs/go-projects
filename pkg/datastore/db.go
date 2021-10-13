package datastore

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"developer.zopsmart.com/go/gofr/pkg"
	"developer.zopsmart.com/go/gofr/pkg/gofr/types"
	"developer.zopsmart.com/go/gofr/pkg/log"
	"developer.zopsmart.com/go/gofr/pkg/middleware"

	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"

	// empty imports are to ensure inits are run for these packages.
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const (
	invalidDialectErr  = "invalid dialect: supported dialects are - mysql, mssql, sqlite, postgres"
	pushMetricDuration = 100
)

type invalidDialect struct{}

func (i invalidDialect) Error() string {
	return invalidDialectErr
}

// DBConfig stores the config variables required to connect to a supported database
type DBConfig struct {
	HostName string
	Username string
	Password string
	Database string
	Port     string
	Dialect  string // supported dialects are - mysql, mssql, sqlite, postgres
	// postgres related config only, accepts disable, allow, prefer, require,
	// verify-ca and verify-full; default is disable
	SSL               string
	ORM               string
	CertificateFile   string
	KeyFile           string
	ConnRetryDuration int
}

type GORMClient struct {
	*gorm.DB
	config *DBConfig
}

type SQLTx struct {
	*sql.Tx
	logger log.Logger
	config *DBConfig
}

type SQLClient struct {
	*sql.DB
	logger log.Logger
	config *DBConfig
}

// nolint:gochecknoglobals // sqlStats has to be a global variable for prometheus
var (
	sqlStats = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "zs_sql_stats",
		Help:    "Histogram for SQL",
		Buckets: []float64{.001, .003, .005, .01, .025, .05, .1, .2, .3, .4, .5, .75, 1, 2, 3, 5, 10, 30},
	}, []string{"type", "host", "database"})

	sqlOpen = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "zs_sql_open_connections",
		Help: "Gauge for sql open connections",
	}, []string{"database", "host"})

	sqlIdle = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "zs_sql_idle_connections",
		Help: "Gauge for sql idle connections",
	}, []string{"database", "host"})

	sqlInUse = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "zs_sql_inUse_connections",
		Help: "Gauge for sql InUse connections",
	}, []string{"database", "host"})

	_ = prometheus.Register(sqlStats)
	_ = prometheus.Register(sqlOpen)
	_ = prometheus.Register(sqlIdle)
	_ = prometheus.Register(sqlInUse)
)

func pushConnMetrics(database, hostname string, db *sql.DB) {
	for {
		stats := db.Stats()
		sqlOpen.WithLabelValues(database, hostname).Set(float64(stats.OpenConnections))
		sqlIdle.WithLabelValues(database, hostname).Set(float64(stats.Idle))
		sqlInUse.WithLabelValues(database, hostname).Set(float64(stats.InUse))
		time.Sleep(pushMetricDuration * time.Millisecond)
	}
}

// NewORM returns a new ORM object if the config is correct, otherwise it returns the error thrown
func NewORM(config *DBConfig) (GORMClient, error) {
	validDialects := map[string]bool{
		"mysql":    true,
		"mssql":    true,
		"postgres": true,
		"sqlite":   true,
	}

	if _, ok := validDialects[config.Dialect]; !ok {
		return GORMClient{config: config}, invalidDialect{}
	}

	connectionStr := formConnectionStr(config)

	db, err := gorm.Open(config.Dialect, connectionStr)
	if err != nil {
		return GORMClient{config: config}, err
	}

	go pushConnMetrics(config.Database, config.HostName, db.DB())

	return GORMClient{DB: db, config: config}, err
}

// NewORMFromEnv fetches the config from environment variables and returns a new ORM object if the config was
// correct, otherwise returns the thrown error
// Deprecated: Instead use datastore.NewORM
func NewORMFromEnv() (GORMClient, error) {
	// pushing deprecated feature count to prometheus
	middleware.PushDeprecatedFeature("NewORMFromEnv")

	return NewORM(dbConfigFromEnv())
}

type SQLXClient struct {
	*sqlx.DB
	config *DBConfig
}

// NewSQLX returns a new sqlx.DB object if the given config is correct, otherwise throws an error
func NewSQLX(config *DBConfig) (SQLXClient, error) {
	connectionStr := formConnectionStr(config)

	DB, err := sqlx.Connect(config.Dialect, connectionStr)
	if err != nil {
		return SQLXClient{config: config}, err
	}

	go pushConnMetrics(config.Database, config.HostName, DB.DB)

	return SQLXClient{DB: DB, config: config}, err
}

// dbConfigFromEnv fetches the DBConfig from environment
func dbConfigFromEnv() *DBConfig {
	return &DBConfig{
		HostName:        os.Getenv("DB_HOST"),
		Username:        os.Getenv("DB_USER"),
		Password:        os.Getenv("DB_PASSWORD"),
		Database:        os.Getenv("DB_NAME"),
		Port:            os.Getenv("DB_PORT"),
		Dialect:         os.Getenv("DB_DIALECT"),
		SSL:             os.Getenv("DB_SSL"),
		CertificateFile: os.Getenv("DB_CERTIFICATE_FILE"),
		KeyFile:         os.Getenv("DB_KEY_FILE"),
	}
}

// formConnection string forms a DB connection string based on the DB Dialect and the given configuration
func formConnectionStr(config *DBConfig) string {
	switch config.Dialect {
	case "postgres":
		ssl := strings.ToLower(config.SSL)
		if ssl == "" {
			config.SSL = "disable"
		}

		return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s sslkey=%s sslcert=%s",
			config.HostName, config.Port, config.Username, config.Database, config.Password, config.SSL, config.KeyFile, config.CertificateFile)
	case "mssql":
		return fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
			config.Username, config.Password, config.HostName, config.Port, config.Database)
	default: // defaults to mysql
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			config.Username, config.Password, config.HostName, config.Port, config.Database)
	}
}

func (c GORMClient) HealthCheck() types.Health {
	resp := types.Health{
		Name:     pkg.SQL,
		Status:   pkg.StatusDown,
		Host:     c.config.HostName,
		Database: c.config.Database,
	}

	// The following check is for the condition when the connection to SQL has not been made during initialization
	if c.DB == nil {
		return resp
	}

	err := c.DB.DB().Ping()
	if err != nil {
		return resp
	}

	resp.Status = pkg.StatusUp

	return resp
}

func (c SQLXClient) HealthCheck() types.Health {
	resp := types.Health{
		Name:     pkg.SQL,
		Status:   pkg.StatusDown,
		Host:     c.config.HostName,
		Database: c.config.Database,
	}
	// The following check is for the condition when the connection to SQLX has not been made during initialization
	if c.DB == nil {
		return resp
	}

	err := c.DB.Ping()
	if err != nil {
		return resp
	}

	resp.Status = pkg.StatusUp

	return resp
}
