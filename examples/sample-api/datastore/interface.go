package datastore

import (
	"developer.zopsmart.com/go/gofr/examples/sample-api/entity"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Employee interface {
	Post(ctx *gofr.Context, emp entity.Employee) (entity.Employee, error)
	Put(ctx *gofr.Context, id string, emp entity.Employee) (entity.Employee, error)
	Delete(ctx *gofr.Context, id string) (int, error)
	Get(ctx *gofr.Context, id string) (entity.Employee, error)
	GetAll(ctx *gofr.Context) ([]entity.Employee, error)
}
