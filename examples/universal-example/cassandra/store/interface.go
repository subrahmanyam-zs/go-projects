package store

import (
	"github.com/zopsmart/gofr/examples/universal-example/cassandra/entity"
	"github.com/zopsmart/gofr/pkg/gofr"
)

type Employee interface {
	Get(ctx *gofr.Context, filter entity.Employee) []entity.Employee
	Create(ctx *gofr.Context, data entity.Employee) ([]entity.Employee, error)
}
