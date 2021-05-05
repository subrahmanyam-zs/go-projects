package store

import (
	"github.com/zopsmart/gofr/examples/using-cassandra/entity"
	"github.com/zopsmart/gofr/pkg/gofr"
)

type Person interface {
	Get(ctx *gofr.Context, filter entity.Person) []*entity.Person
	Create(ctx *gofr.Context, data entity.Person) ([]*entity.Person, error)
	Delete(ctx *gofr.Context, id string) error
	Update(ctx *gofr.Context, data entity.Person) ([]*entity.Person, error)
}
