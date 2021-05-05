package cache

import "time"

type Cacher interface {
	Get(key string) ([]byte, error)
	Set(key string, content []byte, duration time.Duration) error
	Delete(key string) error
}
