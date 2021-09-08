package datastore

import (
	"errors"
	"io"
	"net"
	"testing"

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

		_, err = db.LogMode(true).DB().Exec("SELECT User FROM mysql.user")
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
	const testDB = "test"

	db := new(DBConfig)
	db.Port = "1234"
	db.HostName = "host"
	db.Password = "pass"
	db.Database = testDB
	db.Username = "user"

	type args struct {
		config *DBConfig
	}

	tests := []struct {
		name    string
		dialect string
		args    args
		want    string
	}{
		{
			name:    "postgres",
			dialect: "postgres",
			args:    args{config: db},
			want:    "host=host port=1234 user=user dbname=test password=pass sslmode=disable sslkey= sslcert=",
		},
		{
			name:    "mssql",
			dialect: "mssql",
			args:    args{config: db},
			want:    "sqlserver://user:pass@host:1234?database=test",
		},
		{
			name:    "mysql",
			dialect: "mysql",
			args:    args{config: db},
			want:    "user:pass@tcp(host:1234)/test?charset=utf8&parseTime=True&loc=Local",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.args.config.Dialect = tt.dialect

			if got := formConnectionStr(tt.args.config); got != tt.want {
				t.Errorf("formConnectionStr() = %v, want %v", got, tt.want)
			}
		})
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

		clientSQL.Close()

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
