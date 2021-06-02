package store

import (
	"developer.zopsmart.com/go/gofr/examples/using-ycql/entity"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Shop interface {
	Get(ctx *gofr.Context, filter entity.Shop) []entity.Shop
	Create(ctx *gofr.Context, data entity.Shop) ([]entity.Shop, error)
	Delete(ctx *gofr.Context, id string) error
	Update(ctx *gofr.Context, data entity.Shop) ([]entity.Shop, error)
}
