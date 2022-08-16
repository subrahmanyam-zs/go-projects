package service

import (
	entities2 "EmployeeDepartment/entities"
	"github.com/google/uuid"
)

type Employee interface {
	Create(employee entities2.Employee) (entities2.Employee, error)
	Update(id uuid.UUID, employee entities2.Employee) (entities2.Employee, error)
	Delete(id uuid.UUID) (int, error)
	Read(id uuid.UUID) (entities2.EmployeeAndDepartment, error)
	ReadAll(name string, includeDepartment bool) ([]entities2.EmployeeAndDepartment, error)
}
type Department interface {
	Create(department entities2.Department) (entities2.Department, error)
	Update(id int, departmet entities2.Department) (entities2.Department, error)
	Delete(id int) (int, error)
}
