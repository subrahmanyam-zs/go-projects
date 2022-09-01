package stores

import (
	"developer.zopsmart.com/go/gofr/examples/data-layer-with-mongo/models"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Customer interface {
	Get(ctx *gofr.Context, name string) ([]models.Customer, error)
	Create(ctx *gofr.Context, model models.Customer) error
	Delete(ctx *gofr.Context, name string) (int, error)
}
