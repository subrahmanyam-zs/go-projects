package store

import (
	"time"

	"github.com/zopsmart/gofr/pkg/errors"
	"github.com/zopsmart/gofr/pkg/gofr"
)

// Store is an abstraction for the core layer
type Store interface {
	Get(ctx *gofr.Context, key string) (string, error)
	Set(ctx *gofr.Context, key, value string, expiration time.Duration) error
	Delete(ctx *gofr.Context, key string) error
}

// Model is the type on which all the core layer's functionality is implemented on
type Model struct{}

// New returns a Model core
func New() *Model {
	return &Model{}
}

// Get returns the value for a given key, throws an error, if something goes wrong
func (m Model) Get(c *gofr.Context, key string) (string, error) {
	// fetch the Redis client
	rc := c.Redis

	value, err := rc.Get(c.Context, key).Result()
	if err != nil {
		return "", errors.DB{Err: err}
	}

	return value, nil
}

// Set accepts a key-value pair, and sets those in Redis, if expiration is non-zero value, it sets a expiration(TTL)
// on those keys, if expiration is 0, then the keys have no expiration time
func (m Model) Set(c *gofr.Context, key, value string, expiration time.Duration) error {
	// fetch the Redis client
	rc := c.Redis

	if err := rc.Set(c.Context, key, value, expiration).Err(); err != nil {
		return errors.DB{Err: err}
	}

	return nil
}

// Delete deletes a key from Redis, returns the error if it fails to delete
func (m Model) Delete(c *gofr.Context, key string) error {
	// fetch the Redis client
	rc := c.Redis
	return rc.Del(c.Context, key).Err()
}
