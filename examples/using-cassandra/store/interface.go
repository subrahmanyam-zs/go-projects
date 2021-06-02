package store

import (
	"developer.zopsmart.com/go/gofr/examples/using-cassandra/entity"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Person interface {
	Get(ctx *gofr.Context, filter entity.Person) []*entity.Person
	Create(ctx *gofr.Context, data entity.Person) ([]*entity.Person, error)
	Delete(ctx *gofr.Context, id string) error
	Update(ctx *gofr.Context, data entity.Person) ([]*entity.Person, error)
}
