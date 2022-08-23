package datastore

import (
	"errors"
	"io"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/pkg"
	"developer.zopsmart.com/go/gofr/pkg/gofr/config"
	"developer.zopsmart.com/go/gofr/pkg/log"
)

func TestNewORM(t *testing.T) {
	// failure case
	{
		_, err := NewORM(&DBConfig{
			HostName: "fake host",
			Username: "root",
			Password: "root123",
			Database: "mysql",
			Port:     "1000",
			Dialect:  "mysql",
		})

		e := new(net.DNSError)

		if err != nil && !errors.As(err, &e) {
			t.Errorf("FAILED, expected: %s, got: %s", e, err)
		}
	}

	// failure case due to invalid dialect
	{
		dc := DBConfig{
			Dialect: "fake dialect",
		}

		_, err := NewORM(&dc)
		if err != nil && !errors.As(err, &invalidDialect{}) {
			t.Errorf("FAILED, expected: %v, got: %v", invalidDialect{}.Error(), err)
		}
	}

	// success case
	{
		c := config.NewGoDotEnvProvider(log.NewMockLogger(io.Discard), "../../configs")

		dc := DBConfig{
			HostName: c.Get("DB_HOST"),
			Username: c.Get("DB_USER"),
			Password: c.Get("DB_PASSWORD"),
			Database: c.Get("DB_NAME"),
			Port:     c.Get("DB_PORT"),
			Dialect:  c.Get("DB_DIALECT"),
		}

		db, err := NewORM(&dc)
		if err != nil {
			t.Errorf("FAILED, Could not connect to SQL, got error: %v\n", err)
			return
		}

		err = db.Exec("SELECT User FROM mysql.user").Error
		if err != nil {
			t.Errorf("FAILED, Could not run sql command, got error: %v\n", err)
		}
	}
}

func TestInvalidDialect_Error(t *testing.T) {
	var err invalidDialect

	expected := "invalid dialect: supported dialects are - mysql, mssql, sqlite, postgres"

	if err.Error() != expected {
		t.Errorf("FAILED, Expected: %v, Got: %v", expected, err)
	}
}

func Test_formConnectionStr(t *testing.T) {
	cfg := DBConfig{
		HostName: "host",
		Username: "user",
		Password: "pass",
		Database: "test",
		Port:     "1234",
	}

	tests := []struct {
		name    string
		dialect string
		want    string
	}{
		{"postgres", "postgres", "postgres://user@host:1234/test?password=pass&sslmode=disable&sslcert=&sslkey="},
		{"mssql", "mssql", "sqlserver://user:pass@host:1234?database=test"},
		{"mysql", "mysql", "user:pass@tcp(host:1234)/test?charset=utf8&parseTime=True&loc=Local"},
	}

	for i, tc := range tests {
		cfg.Dialect = tc.dialect

		got := formConnectionStr(&cfg)

		assert.Equal(t, tc.want, got, "TEST[%v] failed\n%s", i, tc.name)
	}
}

func Test_NewSQLX(t *testing.T) {
	// failure case
	{
		_, err := NewSQLX(&DBConfig{
			HostName: "fake host",
			Username: "root",
			Password: "root123",
			Database: "mysql",
			Port:     "1000",
			Dialect:  "mysql",
		})

		e := new(net.DNSError)

		if err != nil && !errors.As(err, &e) {
			t.Errorf("FAILED, expected: %s, got: %s", e, err)
		}
	}

	// success case
	{
		c := config.NewGoDotEnvProvider(log.NewMockLogger(io.Discard), "../../configs")

		dc := DBConfig{
			HostName: c.Get("DB_HOST"),
			Username: c.Get("DB_USER"),
			Password: c.Get("DB_PASSWORD"),
			Database: c.Get("DB_NAME"),
			Port:     c.Get("DB_PORT"),
			Dialect:  c.Get("DB_DIALECT"),
		}

		_, err := NewSQLX(&dc)
		if err != nil {
			t.Errorf("FAILED, expected: %v, got: %v", nil, err)
		}
	}
}

func TestDataStore_SQL_SQLX_HealthCheck(t *testing.T) {
	c := config.NewGoDotEnvProvider(log.NewMockLogger(io.Discard), "../../configs")

	dbConfig := DBConfig{HostName: c.Get("DB_HOST"), Username: c.Get("DB_USER"), Password: c.Get("DB_PASSWORD"),
		Database: c.Get("DB_NAME"), Port: c.Get("DB_PORT"), Dialect: c.Get("DB_DIALECT"),
	}

	testcases := []struct {
		host   string
		status string
	}{
		{dbConfig.HostName, pkg.StatusUp},
		{"invalid", pkg.StatusDown},
	}

	for i, v := range testcases {
		dbConfig.HostName = v.host

		clientSQL, _ := NewORM(&dbConfig)
		dsSQL := DataStore{gorm: clientSQL}

		healthCheck := dsSQL.SQLHealthCheck()
		if healthCheck.Status != v.status {
			t.Errorf("[TESTCASE%d]SQL Failed. Expected status: %v\n Got: %v", i+1, v.status, healthCheck)
		}

		// connecting to SQLX
		clientSQLX, _ := NewSQLX(&dbConfig)
		dsSQLX := DataStore{sqlx: clientSQLX}

		healthCheck = dsSQLX.SQLXHealthCheck()
		if healthCheck.Status != v.status {
			t.Errorf("[TESTCASE%d]SQLX Failed. Expected status: %v\n Got: %v", i+1, v.status, healthCheck)
		}
	}
}

// Test_SQL_SQLX_HealthCheck_Down tests health check response when the db connection was made but lost in between
func Test_SQL_SQLX_HealthCheck_Down(t *testing.T) {
	c := config.NewGoDotEnvProvider(log.NewMockLogger(io.Discard), "../../configs")

	dbConfig := DBConfig{
		HostName: c.Get("DB_HOST"),
		Username: c.Get("DB_USER"),
		Password: c.Get("DB_PASSWORD"),
		Database: c.Get("DB_NAME"),
		Port:     c.Get("DB_PORT"),
		Dialect:  c.Get("DB_DIALECT"),
	}

	{
		clientSQL, _ := NewORM(&dbConfig)
		dsSQL := DataStore{gorm: clientSQL}

		db, _ := clientSQL.DB.DB()

		// db connected but goes down in between
		db.Close()

		healthCheck := dsSQL.SQLHealthCheck()
		if healthCheck.Status != pkg.StatusDown {
			t.Errorf("Failed. Expected: DOWN, Got: %v", healthCheck.Status)
		}
	}

	{
		// connecting to SQLX
		clientSQLX, _ := NewSQLX(&dbConfig)
		dsSQLX := DataStore{sqlx: clientSQLX}

		// db connected but goes down in between
		clientSQLX.Close()

		healthCheck := dsSQLX.SQLXHealthCheck()
		if healthCheck.Status != pkg.StatusDown {
			t.Errorf("Failed. Expected: DOWN, Got: %v", healthCheck.Status)
		}
	}
}
