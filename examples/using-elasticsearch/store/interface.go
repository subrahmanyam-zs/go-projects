package store

import (
	"github.com/zopsmart/gofr/examples/using-elasticsearch/model"
	"github.com/zopsmart/gofr/pkg/gofr"
)

type Customer interface {
	Get(context *gofr.Context, name string) ([]model.Customer, error)
	GetByID(context *gofr.Context, id string) (*model.Customer, error)
	Update(context *gofr.Context, customer model.Customer, id string) (*model.Customer, error)
	Create(context *gofr.Context, customer model.Customer) (*model.Customer, error)
	Delete(context *gofr.Context, id string) error
}
