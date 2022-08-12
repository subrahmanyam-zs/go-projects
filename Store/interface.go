package Store

import (
	"EmployeeDepartment/Handler/Entities"
	"github.com/google/uuid"
)

type Employee interface {
	Create(employee Entities.Employee) (Entities.Employee, error)
	Read(uid uuid.UUID) ([]Entities.EmployeeAndDepartment, error)
	Update(uid string, employee Entities.Employee) (Entities.Employee, error)
	Delete(uid uuid.UUID) (int, error)
	ReadAll(name string, includeDepartment bool) ([]Entities.EmployeeAndDepartment, error)
	ReadDepartment(id int) (Entities.Department, error)
}

type Department interface {
	Create(department Entities.Department) (Entities.Department, error)
	Update(id int, department Entities.Department) (Entities.Department, error)
	Delete(id int) (int, error)
}
