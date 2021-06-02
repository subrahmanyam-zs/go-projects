package cache

import (
	"context"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/datastore"
)

type RedisCacher struct {
	redis datastore.Redis
}

func NewRedisCacher(redis datastore.Redis) RedisCacher {
	return RedisCacher{redis: redis}
}

func (r RedisCacher) Get(key string) ([]byte, error) {
	resp, err := r.redis.Get(context.Background(), key).Bytes()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r RedisCacher) Set(key string, content []byte, duration time.Duration) error {
	return r.redis.Set(context.Background(), key, content, duration).Err()
}

func (r RedisCacher) Delete(key string) error {
	return r.redis.Del(context.Background(), key).Err()
}
