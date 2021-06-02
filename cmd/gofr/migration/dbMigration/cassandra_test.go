package dbmigration

import (
	"strconv"
	"testing"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/gofr/config"
	"developer.zopsmart.com/go/gofr/pkg/log"
)

func TestCassandra_IsDirty(t *testing.T) {
	logger := log.NewLogger()
	c := config.NewGoDotEnvProvider(logger, "../../../../configs")

	// Initialize Cassandra for Dirty Check
	port, _ := strconv.Atoi(c.Get("CASS_DB_PORT"))
	cassandra, _ := datastore.GetNewCassandra(logger, &datastore.CassandraCfg{
		Hosts:    c.Get("CASS_DB_HOST"),
		Port:     port,
		Username: c.Get("CASS_DB_USER"),
		Password: c.Get("CASS_DB_PASS"),
		Keyspace: "test"})

	createCassandraTable(&cassandra, t)

	defer func() {
		err := cassandra.Session.Query("DROP TABLE IF EXISTS gofr_migrations ").Exec()
		if err != nil {
			t.Errorf("Got error while dropping the table at last: %v", err)
		}
	}()

	check := NewCassandra(&cassandra).isDirty("testingCassandra")
	if !check {
		t.Errorf("Dirty migration check for cassandra is false")
	}
}

// createCassandraTable will create a fresh table called gofr_migrations and insert the data required for TestCassandra_IsDirty
func createCassandraTable(cassandra *datastore.Cassandra, t *testing.T) {
	err := cassandra.Session.Query("DROP TABLE IF EXISTS gofr_migrations ").Exec()
	if err != nil {
		t.Errorf("Got error while dropping the existing table gofr_migrations: %v", err)
	}

	migrationTableSchema := "CREATE TABLE IF NOT EXISTS gofr_migrations ( app text, version bigint, start_time timestamp, " +
		"end_time timestamp, method text, PRIMARY KEY (app, version, method) )"
	err = cassandra.Session.Query(migrationTableSchema).Exec()

	if err != nil {
		t.Errorf("Failed creation of gofr_migrations table :%v", err)
	}

	err = cassandra.Session.Query("INSERT INTO gofr_migrations (app, version, start_time, method, end_time) "+
		"values ('testingCassandra', 7, dateof(now()), 'UP', ?)", time.Time{}).Exec()
	if err != nil {
		t.Errorf("Insert Error: %v", err)
	}
}
