package dbmigration

import (
	"context"
	"testing"

	"github.com/zopsmart/gofr/pkg/datastore"
	"github.com/zopsmart/gofr/pkg/errors"
	"github.com/zopsmart/gofr/pkg/log"
)

type K20200324120906 struct {
}

func (k K20200324120906) Up(d *datastore.DataStore, l log.Logger) error {
	l.Info("Running test migration: UP")
	return nil
}

func (k K20200324120906) Down(d *datastore.DataStore, l log.Logger) error {
	return &errors.Response{Reason: "test error"}
}

func TestRedis_Run(t *testing.T) {
	logger := log.NewLogger()
	ctx := context.Background()
	redis, _ := datastore.NewRedisFromEnv(nil)

	redis.Incr(ctx, "redisTest"+migrationLock)

	defer func() {
		redis.Del(ctx, "redisTest"+migrationLock)
	}()

	type args struct {
		mig    Migrator
		app    string
		name   string
		method string
		logger log.Logger
	}

	tt := struct {
		name    string
		args    args
		wantErr bool
	}{"lock acquired", args{K20200324120906{}, "redisTest", "20200324120906", "UP", logger}, false}

	r := NewRedis(redis)
	if err := r.Run(tt.args.mig, tt.args.app, tt.args.name, tt.args.method, tt.args.logger); (err != nil) != tt.wantErr {
		t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
	}
}

func TestRedis_DOWN(t *testing.T) {
	database, _ := datastore.NewRedisFromEnv(nil)

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

	r := NewRedis(database)
	if err := r.Run(K20180324120906{}, tt.args.app, "20180324120906", tt.args.method, log.NewLogger()); (err != nil) != tt.wantErr {
		t.Errorf("postRun() error = %v, wantErr %v", err, tt.wantErr)
	}
}
