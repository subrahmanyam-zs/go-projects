package dbmigration

import (
	"io"
	"strconv"
	"testing"

	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/gofr/config"
	"developer.zopsmart.com/go/gofr/pkg/log"
)

func fetchYCQL(t *testing.T) YCQL {
	logger := log.NewMockLogger(io.Discard)
	c := config.NewGoDotEnvProvider(logger, "../../../../configs")

	port, err := strconv.Atoi(c.Get("YCQL_DB_PORT"))
	if err != nil {
		port = 9042
	}

	yugabyteDBConfig := datastore.CassandraCfg{
		Hosts:       c.Get("CASS_DB_HOST"),
		Port:        port,
		Consistency: datastore.LocalQuorum,
		Username:    c.Get("YCQL_DB_USER"),
		Password:    c.Get("YCQL_DB_PASS"),
		Keyspace:    c.Get("CASS_DB_KEYSPACE"),
		Timeout:     600,
	}

	db, err := datastore.GetNewYCQL(logger, &yugabyteDBConfig)

	if err != nil {
		t.Error(err)
	}

	ycqlDB := NewYCQL(&db)

	return *ycqlDB
}

func Test_ycqlMethods(t *testing.T) {
	ycqlDB := fetchYCQL(t)

	migrationTableSchema := "CREATE TABLE IF NOT EXISTS gofr_migrations ( " +
		"app text, version bigint, start_time timestamp, end_time timestamp, method text, PRIMARY KEY (app, version, method) )"

	_ = ycqlDB.session.Query(migrationTableSchema).Exec()

	if !ycqlDB.isDirty("appName") {
		t.Errorf("Failed")
	}

	testcases := struct {
		app     string
		method  string
		name    string
		wantErr bool
	}{
		"appName", "UP", "K20210116140839", true,
	}

	if err := ycqlDB.preRun(testcases.app, testcases.method, testcases.name); (err != nil) != testcases.wantErr {
		t.Errorf("Failed. Got %s", err)
	}

	if err := ycqlDB.postRun(testcases.app, testcases.method, testcases.name); err != nil {
		t.Errorf("Failed. Got %s", err)
	}

	if err := ycqlDB.FinishMigration(); err != nil {
		t.Errorf("Failed. Got %s", err)
	}
}
