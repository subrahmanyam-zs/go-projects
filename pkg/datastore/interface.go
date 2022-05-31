package datastore

import "context"

type tester interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
}

// KVData implements methods to store,retrieve and delete key value pairs.
type KVData interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string) error
	Delete(ctx context.Context, key string) error
}
