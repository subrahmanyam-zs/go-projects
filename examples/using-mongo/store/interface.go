package store

import (
	"github.com/zopsmart/gofr/examples/using-mongo/entity"
	"github.com/zopsmart/gofr/pkg/gofr"
)

type Customer interface {
	Get(ctx *gofr.Context, name string) ([]*entity.Customer, error)
	Create(ctx *gofr.Context, model *entity.Customer) error
	Delete(ctx *gofr.Context, name string) (int, error)
}
