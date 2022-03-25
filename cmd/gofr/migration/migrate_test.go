package migration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	dbmigration "developer.zopsmart.com/go/gofr/cmd/gofr/migration/dbMigration"
	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr/config"
	"developer.zopsmart.com/go/gofr/pkg/log"
)

type K20200324120906 struct{}

func (k K20200324120906) Up(d *datastore.DataStore, l log.Logger) error {
	l.Info("Running test migration: UP")
	return nil
}

func (k K20200324120906) Down(d *datastore.DataStore, l log.Logger) error {
	return &errors.Response{Reason: "test error"}
}

type K20200324150906 struct{}

func (k K20200324150906) Up(d *datastore.DataStore, l log.Logger) error {
	l.Info("Running test migration: UP")
	return nil
}

func (k K20200324150906) Down(d *datastore.DataStore, l log.Logger) error {
	return &errors.Response{Reason: "test error"}
}

type K20190324150906 struct{}

func (k K20190324150906) Up(d *datastore.DataStore, l log.Logger) error {
	l.Info("Running test migration: UP")
	return nil
}

func (k K20190324150906) Down(d *datastore.DataStore, l log.Logger) error {
	return &errors.Response{Reason: "test error"}
}

type K20200402143245 struct{}

func (k K20200402143245) Up(d *datastore.DataStore, l log.Logger) error {
	l.Info("Running test migration: UP")
	return &errors.Response{Reason: "test error"}
}

func (k K20200402143245) Down(d *datastore.DataStore, l log.Logger) error {
	return nil
}

type K20200423083024 struct{}

func (k K20200423083024) Up(d *datastore.DataStore, logger log.Logger) error {
	return nil
}

func (k K20200423083024) Down(d *datastore.DataStore, logger log.Logger) error {
	return nil
}

type K20200423093024 struct{}

func (k K20200423093024) Up(d *datastore.DataStore, logger log.Logger) error {
	return nil
}

func (k K20200423093024) Down(d *datastore.DataStore, logger log.Logger) error {
	return nil
}

const (
	appName  = "gofr-test"
	keyspace = "cassandra_test"
)

func TestMain(m *testing.M) {
	logger := log.NewLogger()
	c := config.NewGoDotEnvProvider(logger, "../../../configs")
	cassandraPort, _ := strconv.Atoi(c.Get("CASS_DB_PORT"))
	cassandraCfg := datastore.CassandraCfg{
		Hosts:    c.Get("CASS_DB_HOST"),
		Port:     cassandraPort,
		Username: c.Get("CASS_DB_USER"),
		Password: c.Get("CASS_DB_PASS"),
		Keyspace: "system",
	}

	cassandra, err := datastore.GetNewCassandra(logger, &cassandraCfg)
	if err != nil {
		logger.Errorf("[FAILED] unable to connect to cassandra with system keyspace %s", err)
	}

	query := fmt.Sprintf("CREATE KEYSPACE IF NOT EXISTS %v WITH replication = "+
		"{'class':'SimpleStrategy', 'replication_factor' : 1} ", keyspace)

	err = cassandra.Session.Query(query).Exec()
	if err != nil {
		logger.Errorf("unable to create %v keyspace %s", keyspace, err)
	}

	os.Exit(m.Run())
}

func TestRedisAndMongo_Migration(t *testing.T) {
	logger := log.NewMockLogger(new(bytes.Buffer))
	c := config.NewGoDotEnvProvider(log.NewMockLogger(new(bytes.Buffer)), "../../../configs")

	// initialize data stores
	redis, _ := datastore.NewRedis(logger, datastore.RedisConfig{
		HostName: c.Get("REDIS_HOST"),
		Port:     c.Get("REDIS_PORT"),
	})

	mongo, _ := datastore.GetNewMongoDB(logger, &datastore.MongoConfig{
		HostName: c.Get("MONGO_DB_HOST"),
		Port:     c.Get("MONGO_DB_PORT"),
		Username: c.Get("MONGO_DB_USER"),
		Password: c.Get("MONGO_DB_PASS"),
		Database: c.Get("MONGO_DB_NAME")})

	defer func() {
		_ = mongo.Collection("gofr_migrations").Drop(context.TODO())

		redis.Del(context.Background(), "gofr_migrations")
	}()

	testcases := []struct {
		method     string
		migrations map[string]dbmigration.Migrator

		err error
	}{
		{"UP", map[string]dbmigration.Migrator{"20190324150906": K20190324150906{}}, nil},
		{"UP", nil, nil},
		{"UP", map[string]dbmigration.Migrator{"20200324120906": K20200324120906{}}, nil},
		{"UP", map[string]dbmigration.Migrator{"20200402143245": K20200402143245{}}, &errors.Response{Reason: "test error"}},
		{"DOWN", map[string]dbmigration.Migrator{"20200324120906": K20200324120906{}}, &errors.Response{Reason: "test error"}},
		{"UP", map[string]dbmigration.Migrator{"20200423083024": K20200423083024{}, "20200423093024": K20200423093024{}}, nil},
		{"DOWN", map[string]dbmigration.Migrator{"20200423083024": K20200423083024{}, "20200423093024": K20200423093024{}}, nil},
		{"DOWN", map[string]dbmigration.Migrator{"20200423083024": K20200423083024{}}, nil},
	}

	for i, v := range testcases {
		err := Migrate(appName, dbmigration.NewRedis(redis), v.migrations, v.method, logger)
		if !reflect.DeepEqual(err, v.err) {
			t.Errorf("[TESTCASE%d]Redis : Got %v\tExpected %v\n", i+1, err, v.err)
		}

		err = Migrate(appName, dbmigration.NewMongo(mongo), v.migrations, v.method, logger)
		if !reflect.DeepEqual(err, v.err) {
			t.Errorf("[TESTCASE%d]Mongo : Got %v\tExpected %v\n", i+1, err, v.err)
		}
	}
}

func TestMySQL_Migration(t *testing.T) {
	logger := log.NewMockLogger(new(bytes.Buffer))
	c := config.NewGoDotEnvProvider(log.NewMockLogger(new(bytes.Buffer)), "../../../configs")
	mysql, _ := datastore.NewORM(&datastore.DBConfig{
		HostName: c.Get("DB_HOST"),
		Username: c.Get("DB_USER"),
		Password: c.Get("DB_PASSWORD"),
		Database: c.Get("DB_NAME"),
		Port:     c.Get("DB_PORT"),
		Dialect:  "mysql",
	})

	defer func() {
		_ = mysql.Migrator().DropTable("gofr_migrations")
	}()

	// ensures the gofr_migrations table is dropped in DB
	tx := mysql.DB.Exec("DROP TABLE IF EXISTS gofr_migrations")
	if tx != nil {
		assert.NoError(t, tx.Error)
	}

	testcases := []struct {
		method     string
		migrations map[string]dbmigration.Migrator

		err error
	}{
		{"UP", map[string]dbmigration.Migrator{"20190324150906": K20190324150906{}}, nil},
		{"UP", nil, nil},
		{"UP", map[string]dbmigration.Migrator{"20200402143245": K20200402143245{}},
			&errors.Response{Reason: "error encountered in running the migration", Detail: &errors.Response{Reason: "test error"}}},
		{"UP", map[string]dbmigration.Migrator{"20200324120906": K20200324120906{}}, nil},
		{"DOWN", map[string]dbmigration.Migrator{"20200324120906": K20200324120906{}},
			&errors.Response{Reason: "error encountered in running the migration", Detail: &errors.Response{Reason: "test error"}}},
		{"UP", map[string]dbmigration.Migrator{"20200423083024": K20200423083024{}, "20200423093024": K20200423093024{}}, nil},
		{"DOWN", map[string]dbmigration.Migrator{"20200423083024": K20200423083024{}, "20200423093024": K20200423093024{}}, nil},
		{"DOWN", map[string]dbmigration.Migrator{"20200423083024": K20200423083024{}}, nil},
	}

	for i, v := range testcases {
		err := Migrate(appName, dbmigration.NewGorm(mysql.DB), v.migrations, v.method, logger)
		if !reflect.DeepEqual(err, v.err) {
			t.Errorf("[TESTCASE%d]Got %v\tExpected %v\n", i+1, err, v.err)
		}
	}
}

func TestCassandra_Migration(t *testing.T) {
	logger := log.NewMockLogger(new(bytes.Buffer))
	c := config.NewGoDotEnvProvider(logger, "../../../configs")
	cassandraPort, _ := strconv.Atoi(c.Get("CASS_DB_PORT"))
	cassandraCfg := datastore.CassandraCfg{
		Hosts:    c.Get("CASS_DB_HOST"),
		Port:     cassandraPort,
		Username: c.Get("CASS_DB_USER"),
		Password: c.Get("CASS_DB_PASS"),
		Keyspace: keyspace,
	}

	cassandra, _ := datastore.GetNewCassandra(logger, &cassandraCfg)
	_ = cassandra.Session.Query("DROP TABLE IF EXISTS gofr_migrations  ").Exec()

	defer func() {
		_ = cassandra.Session.Query("DROP TABLE IF EXISTS gofr_migrations ").Exec()
	}()

	testcases := []struct {
		method     string
		migrations map[string]dbmigration.Migrator

		err error
	}{
		{"UP", map[string]dbmigration.Migrator{"20190324150906": K20190324150906{}}, nil},
		{"UP", map[string]dbmigration.Migrator{"20200324120906": K20200324120906{}}, nil},
		{"UP", map[string]dbmigration.Migrator{"20200402143245": K20200402143245{}},
			&errors.Response{Reason: "error encountered in running the migration", Detail: &errors.Response{Reason: "test error"}}},
		{"DOWN", map[string]dbmigration.Migrator{"20200324120906": K20200324120906{}},
			&errors.Response{Reason: "error encountered in running the migration", Detail: &errors.Response{Reason: "test error"}}},
		{"UP", map[string]dbmigration.Migrator{"20200423083024": K20200423083024{}, "20200423093024": K20200423093024{}}, nil},
		{"DOWN", map[string]dbmigration.Migrator{"20200423083024": K20200423083024{}, "20200423093024": K20200423093024{}}, nil},
		{"DOWN", map[string]dbmigration.Migrator{"20200423083024": K20200423083024{}}, nil},
	}

	for i, v := range testcases {
		err := Migrate(appName, dbmigration.NewCassandra(&cassandra), v.migrations, v.method, logger)
		if !reflect.DeepEqual(err, v.err) {
			t.Errorf("[TESTCASE%d]Got %v\tExpected %v\n", i+1, err, v.err)
		}
	}
}

func TestMigrateError(t *testing.T) {
	logger := log.NewMockLogger(new(bytes.Buffer))

	err := Migrate(appName, nil, nil, "UP", logger)
	if err == nil {
		t.Errorf("expected err, got nil")
	}
}

func Test_MigrateCheck(t *testing.T) {
	b := new(bytes.Buffer)
	mockLogger := log.NewMockLogger(b)
	c := config.NewGoDotEnvProvider(mockLogger, "../../../configs")

	mysql, _ := datastore.NewORM(&datastore.DBConfig{
		HostName: c.Get("DB_HOST"),
		Username: c.Get("DB_USER"),
		Password: c.Get("DB_PASSWORD"),
		Database: c.Get("DB_NAME"),
		Port:     c.Get("DB_PORT"),
		Dialect:  "mysql",
	})

	defer func() {
		_ = mysql.Migrator().DropTable("gofr_migrations")
	}()

	// ensures the gofr_migrations table is dropped in DB
	tx := mysql.DB.Exec("DROP TABLE IF EXISTS gofr_migrations")
	if tx != nil {
		assert.NoError(t, tx.Error)
	}

	migrations := map[string]dbmigration.Migrator{"20210324150906": K20200324150906{},
		"20200324120906": K20200324120906{},
		"20190324150906": K20190324150906{}}

	if err := Migrate(appName, dbmigration.NewGorm(mysql.DB), migrations, "UP", mockLogger); err != nil {
		t.Errorf("expected nil error\tgot %v", err)
	}

	loggedStr := b.String()
	i1 := strings.Index(loggedStr, "20190324150906")
	i2 := strings.Index(loggedStr, "20200324120906")
	i3 := strings.Index(loggedStr, "20210324150906")

	if i1 > i2 || i2 > i3 {
		t.Errorf("Sequence of migration run is not in order, got: %v", loggedStr)
	}
}

type gofrMigration struct {
	App       string    `gorm:"primary_key"`
	Version   int64     `gorm:"primary_key;auto_increment:false"`
	StartTime time.Time `gorm:"autoCreateTime"`
	EndTime   time.Time `gorm:"default:NULL"`
	Method    string    `gorm:"primary_key"`
}

func Test_DirtyTest(t *testing.T) {
	logger := log.NewMockLogger(new(bytes.Buffer))
	c := config.NewGoDotEnvProvider(log.NewMockLogger(new(bytes.Buffer)), "../../../configs")

	// initialize data stores
	redis, _ := datastore.NewRedis(logger, datastore.RedisConfig{
		HostName: c.Get("REDIS_HOST"),
		Port:     c.Get("REDIS_PORT"),
	})

	port, _ := strconv.Atoi(c.Get("CASS_DB_PORT"))
	cassandra, _ := datastore.GetNewCassandra(logger, &datastore.CassandraCfg{
		Hosts:       c.Get("CASS_DB_HOST"),
		Port:        port,
		Consistency: c.Get("CASS_DB_CONSISTENCY"),
		Username:    c.Get("CASS_DB_CONSISTENCY"),
		Password:    c.Get("CASS_DB_PASS"),
		Keyspace:    "test"})

	_ = cassandra.Session.Query("drop table gofr_migrations").Exec()

	mysql, _ := datastore.NewORM(&datastore.DBConfig{
		HostName: c.Get("DB_HOST"),
		Username: c.Get("DB_USER"),
		Password: c.Get("DB_PASSWORD"),
		Database: c.Get("DB_NAME"),
		Port:     c.Get("DB_PORT"),
		Dialect:  "mysql",
	})

	migrationTableSchema := "CREATE TABLE IF NOT EXISTS gofr_migrations ( app text, version bigint, start_time timestamp, end_time text, " +
		"method text, PRIMARY KEY (app, version, method) )"

	_ = cassandra.Session.Query(migrationTableSchema).Exec()
	_ = cassandra.Session.Query("insert into gofr_migrations (app, version, start_time, method, end_time) " +
		"values ('testing', 12, dateof(now()), 'UP', '')").Exec()

	err := mysql.Migrator().CreateTable(&gofrMigration{})
	assert.NoError(t, err)

	mysql.Create(&gofrMigration{App: "testing", Method: "UP", Version: 20000102121212, StartTime: time.Now()})

	ctx := context.Background()
	redisMigrator := dbmigration.NewRedis(redis)
	resBytes, _ := json.Marshal([]gofrMigration{{"testing", 20010102121212, time.Now(), time.Time{}, "UP"},
		{"testing", 20000102121212, time.Now(), time.Time{}, "UP"}})
	redis.HSet(ctx, "gofr_migrations", "testing", string(resBytes))

	redisMigrator.LastRunVersion("testing", "UP")

	defer func() {
		_ = mysql.Migrator().DropTable("gofr_migrations")
		_ = cassandra.Session.Query("truncate gofr_migrations").Exec()
		_ = redis.Del(ctx, "gofr_migrations")
	}()

	type args struct {
		app    string
		method string
		name   string
	}

	tests := []struct {
		name     string
		dbDriver dbmigration.DBDriver
		args     args
		wantErr  bool
	}{
		{"cassandra: dirty check", dbmigration.NewCassandra(&cassandra), args{"testing", "UP", "12"}, true},
		{"mysql: dirty check", dbmigration.NewGorm(mysql.DB), args{"testing", "UP", "12"}, true},
		{"redis: dirty check", redisMigrator, args{"testing", "UP", "12"}, true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.dbDriver.Run(nil, tt.args.app, tt.args.method, tt.args.name, logger); (err != nil) != tt.wantErr {
				t.Errorf("preRun() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
