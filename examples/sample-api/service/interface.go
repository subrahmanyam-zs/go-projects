package service

import (
	"developer.zopsmart.com/go/gofr/examples/sample-api/entity"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Employee interface {
	Post(ctx *gofr.Context, employee entity.Employee) (interface{}, error)
	Put(ctx *gofr.Context, id string, emp entity.Employee) (interface{}, error)
	Delete(ctx *gofr.Context, id string) (interface{}, error)
	Get(ctx *gofr.Context, id string) (interface{}, error)
	GetAll(ctx *gofr.Context) (interface{}, error)
}
