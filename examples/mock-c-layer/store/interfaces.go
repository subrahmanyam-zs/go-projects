package store

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Brand interface {
	Get(ctx *gofr.Context) ([]Model, error)
	Create(ctx *gofr.Context, brand Model) (Model, error)
}
