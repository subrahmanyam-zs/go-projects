package Store

import "EmployeeDepartment/Handler/Entities"

type Employee interface {
	Create(employee Entities.Employee) (Entities.Employee, error)
}

type Department interface {
	Create(department Entities.Department) (Entities.Department, error)
}
