package store

import (
	entities2 "EmployeeDepartment/entities"
	"github.com/google/uuid"
)

type Employee interface {
	Create(employee entities2.Employee) (entities2.Employee, error)
	Read(uid uuid.UUID) ([]entities2.EmployeeAndDepartment, error)
	Update(uid string, employee entities2.Employee) (entities2.Employee, error)
	Delete(uid uuid.UUID) (int, error)
	ReadAll(name string, includeDepartment bool) ([]entities2.EmployeeAndDepartment, error)
	ReadDepartment(id int) (entities2.Department, error)
}

type Department interface {
	Create(department entities2.Department) (entities2.Department, error)
	Update(id int, department entities2.Department) (entities2.Department, error)
	Delete(id int) (int, error)
}
