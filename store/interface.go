package store

import (
	entities2 "EmployeeDepartment/entities"
	"context"

	"github.com/google/uuid"
)

type Parameters struct {
	Ctx               context.Context
	Name              string
	IncludeDepartment bool
}

type Employee interface {
	Create(ctx context.Context, employee *entities2.Employee) (*entities2.Employee, error)
	Read(ctx context.Context, uid uuid.UUID) (entities2.EmployeeAndDepartment, error)
	Update(ctx context.Context, uid uuid.UUID, employee *entities2.Employee) (*entities2.Employee, error)
	Delete(ctx context.Context, uid uuid.UUID) (int, error)
	ReadAll(para Parameters) ([]entities2.EmployeeAndDepartment, error)
	ReadDepartment(ctx context.Context, id int) (entities2.Department, error)
}

type Department interface {
	Create(ctx context.Context, department entities2.Department) (entities2.Department, error)
	Update(ctx context.Context, id int, department entities2.Department) (entities2.Department, error)
	Delete(ctx context.Context, id int) (int, error)
	GetDepartment(ctx context.Context, id int) (entities2.Department, error)
}
