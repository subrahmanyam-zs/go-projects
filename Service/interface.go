package Service

import (
	"EmployeeDepartment/Handler/Entities"
	"github.com/google/uuid"
)

type Employee interface {
	Create(employee Entities.Employee) (Entities.Employee, error)
	Update(id uuid.UUID, employee Entities.Employee) (Entities.Employee, error)
	Delete(id uuid.UUID) (int, error)
	Read(id uuid.UUID) (Entities.Employee, error)
	ReadAll(name string, includeDepartment bool) (Entities.EmployeeAndDepartment, error)
}
type Department interface {
	Create(department Entities.Department) (Entities.Department, error)
	Update(id int, departmet Entities.Department) (Entities.Department, error)
	Delete(id int) (int, error)
}
