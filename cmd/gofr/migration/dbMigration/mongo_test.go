package dbmigration

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/zopsmart/gofr/pkg/datastore"
	"github.com/zopsmart/gofr/pkg/log"
)

func TestMongo_IsDirty(t *testing.T) {
	logger := log.NewMockLogger(new(bytes.Buffer))

	mongo, _ := datastore.GetMongoDBFromEnv(logger)
	md := NewMongo(mongo)

	defer func() {
		_ = md.coll.Drop(context.TODO())
	}()

	_, _ = md.coll.InsertOne(context.TODO(), gofrMigration{App: "testing", Version: 20170101100101, Method: UP, StartTime: time.Now()})

	type args struct {
		m      Migrator
		app    string
		name   string
		method string
		logger log.Logger
	}

	tt := struct {
		name    string
		args    args
		wantErr bool
	}{"migration UP", args{nil, "testing", "20200324162754", "UP", logger}, true}

	if err := md.Run(tt.args.m, tt.args.app, tt.args.name, tt.args.method, tt.args.logger); (err != nil) != tt.wantErr {
		t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
	}
}

func TestMongo_DOWN(t *testing.T) {
	logger := log.NewMockLogger(new(bytes.Buffer))

	database, _ := datastore.GetMongoDBFromEnv(logger)

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

	m := NewMongo(database)
	if err := m.Run(K20180324120906{}, tt.args.app, "20180324120906", tt.args.method, logger); (err != nil) != tt.wantErr {
		t.Errorf("postRun() error = %v, wantErr %v", err, tt.wantErr)
	}
}
