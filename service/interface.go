package service

import (
	entities2 "EmployeeDepartment/entities"
	"EmployeeDepartment/store"
	"context"

	"github.com/google/uuid"
)

type Employee interface {
	Create(ctx context.Context, employee *entities2.Employee) (*entities2.Employee, error)
	Update(ctx context.Context, id uuid.UUID, employee *entities2.Employee) (*entities2.Employee, error)
	Delete(ctx context.Context, id uuid.UUID) (int, error)
	Read(ctx context.Context, id uuid.UUID) (entities2.EmployeeAndDepartment, error)
	ReadAll(para store.Parameters) ([]entities2.EmployeeAndDepartment, error)
}
type Department interface {
	Create(ctx context.Context, department entities2.Department) (entities2.Department, error)
	Update(ctx context.Context, id int, department entities2.Department) (entities2.Department, error)
	Delete(ctx context.Context, id int) (int, error)
	GetDepartment(ctx context.Context, id int) (entities2.Department, error)
}
