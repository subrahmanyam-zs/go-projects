package store

import (
	"github.com/zopsmart/gofr/examples/using-ycql/entity"
	"github.com/zopsmart/gofr/pkg/gofr"
)

type Shop interface {
	Get(ctx *gofr.Context, filter entity.Shop) []entity.Shop
	Create(ctx *gofr.Context, data entity.Shop) ([]entity.Shop, error)
	Delete(ctx *gofr.Context, id string) error
	Update(ctx *gofr.Context, data entity.Shop) ([]entity.Shop, error)
}
