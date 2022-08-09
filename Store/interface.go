package Store

import (
	"EmployeeDepartment/Handler/Entities"
	"github.com/google/uuid"
)

type Employee interface {
	Create(employee Entities.Employee) (Entities.Employee, error)
	Read(uid uuid.UUID) (Entities.Employee, error)
	Update(uid uuid.UUID, employee Entities.Employee) (Entities.Employee, error)
}

type Department interface {
	Create(department Entities.Department) (Entities.Department, error)
}
