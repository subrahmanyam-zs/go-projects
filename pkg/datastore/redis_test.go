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

	"developer.zopsmart.com/go/gofr/pkg"
	"developer.zopsmart.com/go/gofr/pkg/gofr/config"
	"developer.zopsmart.com/go/gofr/pkg/gofr/types"
	"developer.zopsmart.com/go/gofr/pkg/log"
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

		_ = r.Close()
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

		_ = r.Close()
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

		_ = r.Close()
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

	{
		cfg := RedisConfig{HostName: c.Get("REDIS_HOST"), Port: c.Get("REDIS_PORT")}
		expHealth := types.Health{Name: pkg.Redis, Status: "UP", Host: c.Get("REDIS_HOST")}

		conn, _ := NewRedis(logger, cfg)
		out := conn.HealthCheck()

		if expHealth.Status != out.Status || expHealth.Host != out.Host || expHealth.Name != out.Name {
			t.Errorf("Success case: Expected %v, got %v", expHealth, out)
		}

		if out.Details == nil {
			t.Errorf("success case: details should not be nil")
		} else if _, ok := out.Details.(map[string]map[string]string); !ok {
			t.Errorf("Success case: type of details mismatched got %T", out.Details)
		}
	}

	{
		cfg := RedisConfig{HostName: "Random", Port: c.Get("REDIS_PORT")}
		exp := types.Health{Name: pkg.Redis, Status: "DOWN", Host: "Random"}

		conn, _ := NewRedis(logger, cfg)

		output := conn.HealthCheck()

		if output != exp {
			t.Errorf("Failed case,Got %v Exepcted %v", output, exp)
		}
	}
}

// connection is made and closed later for HealthCheck
func Test_RedisHealthCheckConnClose(t *testing.T) {
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

func Test_getInfoInMap(t *testing.T) {
	logger := log.NewLogger()
	c := config.NewGoDotEnvProvider(logger, "../../configs")

	cfg := RedisConfig{
		HostName:                c.Get("REDIS_HOST"),
		Port:                    c.Get("REDIS_PORT"),
		ConnectionRetryDuration: 2,
		Options:                 new(redis.Options),
	}

	cfg.Options.Addr = cfg.HostName + ":" + cfg.Port

	client := redisClient{
		Client: redis.NewClient(cfg.Options),
		config: cfg,
	}

	out := client.getInfoInMap()
	if out == nil {
		t.Errorf("Info about client connection should not be nil")
	}
}
