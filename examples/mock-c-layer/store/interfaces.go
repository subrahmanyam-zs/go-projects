package store

import (
	"github.com/zopsmart/gofr/pkg/gofr"
)

type Brand interface {
	Get(ctx *gofr.Context) ([]Model, error)
	Create(ctx *gofr.Context, brand Model) (Model, error)
}
