package store

import (
	"developer.zopsmart.com/go/gofr/examples/using-mongo/entity"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Customer interface {
	Get(ctx *gofr.Context, name string) ([]*entity.Customer, error)
	Create(ctx *gofr.Context, model *entity.Customer) error
	Delete(ctx *gofr.Context, name string) (int, error)
}
