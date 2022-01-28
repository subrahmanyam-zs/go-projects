package store

import (
	"developer.zopsmart.com/go/gofr/examples/mock-c-layer/models"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Brand interface {
	Get(ctx *gofr.Context) ([]models.Brand, error)
	Create(ctx *gofr.Context, brand models.Brand) (models.Brand, error)
}
