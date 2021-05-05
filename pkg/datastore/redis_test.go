package datastore

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/zopsmart/gofr/pkg"
	"github.com/zopsmart/gofr/pkg/gofr/config"
	"github.com/zopsmart/gofr/pkg/gofr/types"
	"github.com/zopsmart/gofr/pkg/log"
)

func Test_NewRedis(t *testing.T) {
	logger := log.NewMockLogger(io.Discard)
	config.NewGoDotEnvProvider(logger, "../../configs")

	{
		// error case
		e := new(net.DNSError)
		e.Err = "address tcp/fake port: unknown port"
		e.Name = "dial tcp"

		if _, err := NewRedis(logger, RedisConfig{
			HostName: "fake host",
			Port:     "6378",
		}); err != nil && !errors.As(err, &e) {
			t.Errorf("FAILED, expected: %s, got: %s", e, err)
		}
	}

	{
		// success case without options
		r, err := NewRedisFromEnv(nil)
		if err != nil {
			t.Error("FAILED, could not connect to Redis: ", err)
			return
		}

		r.Close()
	}

	{
		// success case with options
		r, err := NewRedisFromEnv(&redis.Options{
			MaxRetries: 3,
		})
		if err != nil {
			t.Error("FAILED, could not connect to Redis: ", err)
			return
		}

		r.Close()
	}

	{
		// success case with options, but the Addr is from config.HostName and config.Port
		r, err := NewRedisFromEnv(&redis.Options{
			PoolSize: 5,
		})
		if err != nil {
			t.Error("FAILED, could not connect to Redis: ", err)
			return
		}

		r.Close()
	}
}

func TestNewRedisCluster(t *testing.T) {
	type args struct {
		logger         log.Logger
		clusterOptions *redis.ClusterOptions
	}

	tests := []struct {
		name    string
		args    args
		want    Redis
		wantErr bool
	}{
		{"Error case", args{log.NewLogger(), &redis.ClusterOptions{}}, nil, true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRedisCluster(tt.args.clusterOptions)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewRedisCluster() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRedisCluster() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_RedisQueryLog(t *testing.T) {
	b := new(bytes.Buffer)
	logger := log.NewMockLogger(b)
	c := config.NewGoDotEnvProvider(logger, "../../configs")

	rc := RedisConfig{
		HostName: c.Get("REDIS_HOST"),
		Port:     c.Get("REDIS_PORT"),
	}

	redisClient, _ := NewRedis(logger, rc)

	{ // test query logs
		b.Reset()
		ctx := context.Background()
		_, _ = redisClient.Get(ctx, "someKey123").Result()

		expectedLog := `"get someKey123"`

		if !strings.Contains(b.String(), expectedLog) {
			t.Errorf("[FAILED] expected: %v, got: %v", expectedLog, b.String())
		}

		if !strings.Contains(b.String(), "redis") {
			t.Errorf("[FAILED] expected: %v, got: %v", "REDIS", b.String())
		}
	}

	{ // test batch execution logs
		b.Reset()
		ctx := context.Background()
		_, _ = redisClient.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			_, _ = pipe.Get(ctx, "get-some-key").Result()
			_, _ = pipe.Get(ctx, "someKey123").Result()
			return nil
		})
		expectedLog := `"get get-some-key","get someKey123"`
		if !strings.Contains(b.String(), expectedLog) {
			t.Errorf("[FAILED] expected: %v, got: %v", expectedLog, b.String())
		}

		if !strings.Contains(b.String(), "redis") {
			t.Errorf("[FAILED] expected: %v, got: %v", "REDIS", b.String())
		}
	}
}

func TestDataStore_RedisHealthCheck(t *testing.T) {
	logger := log.NewLogger()
	c := config.NewGoDotEnvProvider(logger, "../../configs")
	testCases := []struct {
		c        RedisConfig
		expected types.Health
	}{
		{RedisConfig{HostName: c.Get("REDIS_HOST"), Port: c.Get("REDIS_PORT")},
			types.Health{Name: pkg.Redis, Status: "UP", Host: c.Get("REDIS_HOST")}},
		{RedisConfig{HostName: "Random", Port: c.Get("REDIS_PORT")},
			types.Health{Name: pkg.Redis, Status: "DOWN", Host: "Random"}},
	}

	for i, tc := range testCases {
		conn, _ := NewRedis(logger, tc.c)
		output := conn.HealthCheck()

		if output != tc.expected {
			t.Errorf("[Failed]%v, Got %v Exepcted %v", i, output, tc.expected)
		}
	}
}

// connection is made and closed later for HealthCheck
func Test_RedisHealthCheck(t *testing.T) {
	logger := log.NewLogger()
	c := config.NewGoDotEnvProvider(logger, "../../configs")
	conf := RedisConfig{HostName: c.Get("REDIS_HOST"), Port: c.Get("REDIS_PORT")}
	expected := types.Health{
		Name: pkg.Redis, Status: "DOWN", Host: c.Get("REDIS_HOST"),
	}

	conn, _ := NewRedis(logger, conf)
	conn.Close()
	output := conn.HealthCheck()

	if output != expected {
		t.Errorf("[Failed] Got %v Exepcted %v", output, expected)
	}
}

// Test for Go-routine leak when redis connection is not established
func Test_goroutineCount(t *testing.T) {
	logger := log.NewLogger()
	c := config.NewGoDotEnvProvider(logger, "../../configs")
	conf := RedisConfig{HostName: c.Get("REDIS_HOST"), Port: "3444"}

	_, _ = NewRedis(logger, conf)
	prev := runtime.NumGoroutine()

	_, _ = NewRedis(logger, conf)
	next := runtime.NumGoroutine()

	if prev != next {
		t.Errorf("[FAILED] Goroutine leaked,Expected: %v,Got: %v", prev, next)
	}
}
