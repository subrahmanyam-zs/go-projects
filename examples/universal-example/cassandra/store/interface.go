package store

import (
	"developer.zopsmart.com/go/gofr/examples/universal-example/cassandra/entity"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Employee interface {
	Get(ctx *gofr.Context, filter entity.Employee) []entity.Employee
	Create(ctx *gofr.Context, data entity.Employee) ([]entity.Employee, error)
}
