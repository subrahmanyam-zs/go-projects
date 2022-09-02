package datastore

import (
	"github.com/google/uuid"

	"developer.zopsmart.com/go/gofr/Emp-Dept/entities"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Employee interface {
	Post(ctx *gofr.Context, emp entities.Employee) (interface{}, error)
	Put(ctx *gofr.Context, id uuid.UUID, emp entities.Employee) (interface{}, error)
	Delete(ctx *gofr.Context, id uuid.UUID) (int, error)
	Get(ctx *gofr.Context, id uuid.UUID) (interface{}, error)
	GetAll(ctx *gofr.Context, name string, includeDept bool) (interface{}, error)
	GetDepartment(ctx *gofr.Context, id int) (entities.Department, error)
}

type Department interface {
	Post(ctx *gofr.Context, dept entities.Department) (interface{}, error)
	Put(ctx *gofr.Context, id int, dept entities.Department) (interface{}, error)
	Delete(ctx *gofr.Context, id int) (int, error)
}
