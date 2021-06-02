package dbmigration

import (
	"io"
	"testing"

	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr/config"
	"developer.zopsmart.com/go/gofr/pkg/log"
)

type K20180324120906 struct {
}

func (k K20180324120906) Up(d *datastore.DataStore, l log.Logger) error {
	query := `
 	   CREATE TABLE IF NOT EXISTS customers (
	   id serial primary key,
	   name varchar (50))
	`
	// nolint:gocritic,sqlclosecheck // returned value not needed
	_, err := d.DB().Query(query)
	if err != nil {
		l.Error("Customer table is not created")
	}

	l.Info("Running test migration: UP")

	return nil
}

func (k K20180324120906) Down(d *datastore.DataStore, l log.Logger) error {
	return &errors.Response{Reason: "test error"}
}

func TestGORM(t *testing.T) {
	c := config.NewGoDotEnvProvider(log.NewMockLogger(io.Discard), "../../../../configs")

	database, _ := datastore.NewORM(&datastore.DBConfig{
		HostName: c.Get("DB_HOST"),
		Username: c.Get("DB_USER"),
		Password: c.Get("DB_PASSWORD"),
		Database: "",
		Port:     c.Get("DB_PORT"),
		Dialect:  c.Get("DB_DIALECT"),
	})
	txn := database.Begin()

	type args struct {
		app    string
		method string
		name   string
	}

	tt := struct {
		name    string
		args    args
		wantErr bool
	}{"database name not selected error", args{"testing", UP, "20200302020202"}, true}

	g := &GORM{
		database: database.DB,
		txn:      txn,
	}
	if err := g.postRun(tt.args.app, tt.args.method, tt.args.name); (err != nil) != tt.wantErr {
		t.Errorf("postRun() error = %v, wantErr %v", err, tt.wantErr)
	}

	if err := g.preRun(tt.args.app, tt.args.method, tt.args.name); (err != nil) != tt.wantErr {
		t.Errorf("preRun() error = %v, wantErr %v", err, tt.wantErr)
	}
}

func TestGORM_DOWN(t *testing.T) {
	database, _ := datastore.NewORMFromEnv()
	txn := database.Begin()

	type args struct {
		app    string
		method string
		ver    int
	}

	tt := struct {
		name    string
		args    args
		wantErr bool
	}{"down error", args{"testing", "DOWN", 20180324120906}, true}

	g := &GORM{
		database: database.DB,
		txn:      txn,
	}
	if err := g.Run(K20180324120906{}, tt.args.app, "20180324120906", tt.args.method, log.NewLogger()); (err != nil) != tt.wantErr {
		t.Errorf("postRun() error = %v, wantErr %v", err, tt.wantErr)
	}
}

func TestGORM_UP(t *testing.T) {
	c := config.NewGoDotEnvProvider(log.NewMockLogger(io.Discard), "../../../../configs")

	database, _ := datastore.NewORM(&datastore.DBConfig{
		HostName: c.Get("DB_HOST"),
		Username: c.Get("DB_USER"),
		Password: c.Get("DB_PASSWORD"),
		Database: "mysql",
		Port:     c.Get("DB_PORT"),
		Dialect:  c.Get("DB_DIALECT"),
	})
	txn := database.Begin()

	g := &GORM{
		database: database.DB,
		txn:      txn,
	}
	if err := g.Run(K20180324120906{}, "testing", "20180324120906", "UP", log.NewLogger()); err != nil {
		t.Errorf("postRun() expected nil error\t got %v", err)
	}
}

func TestGORM_FinishMigration(t *testing.T) {
	config.NewGoDotEnvProvider(log.NewMockLogger(io.Discard), "../../../../configs")

	database, _ := datastore.NewORMFromEnv()
	txn := database.Begin()

	g := &GORM{
		database: database.DB,
		txn:      txn,
	}

	if err := g.FinishMigration(); err != nil {
		t.Errorf("Expected = %v, Got %v", nil, err)
	}
}
