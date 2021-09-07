package store

import (
	"developer.zopsmart.com/go/gofr/examples/using-dynamodb/model"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Person interface {
	Get(c *gofr.Context, id string) (model.Person, error)
	Create(c *gofr.Context, user model.Person) error
	Update(c *gofr.Context, user model.Person) error
	Delete(c *gofr.Context, id string) error
}
