package Department

import (
	"EmployeeDepartment/Handler/Entities"
	"EmployeeDepartment/Store"
	"golang.org/x/exp/slices"
	"strings"
)

type DepartmentHandelr struct {
	dataStore Store.Department
}

func New(dept Store.Department) DepartmentHandelr {
	return DepartmentHandelr{dataStore: dept}
}

func validateName(name string) bool {
	name = strings.ToUpper(name)
	deptNames := []string{"HR", "TECH", "ACCOUNTS"}
	if slices.Contains(deptNames, name) {
		return true
	}
	return false
}

func validateFloorNo(floorNo int) bool {
	if floorNo >= 1 && floorNo <= 3 {
		return true
	}
	return false
}

func (d DepartmentHandelr) validatePost(department Entities.Department) Entities.Department {
	if validateName(department.Name) && validateFloorNo(department.FloorNo) {
		res, err := d.dataStore.Create(department)
		if err != nil {
			return res
		}
		return res
	}
	return Entities.Department{}
}

func (d DepartmentHandelr) validatePut(id int, department Entities.Department) Entities.Department {
	if validateName(department.Name) && validateFloorNo(department.FloorNo) {
		res, err := d.dataStore.Update(id, department)
		if err != nil {
			return res
		}
		return res
	}
	return Entities.Department{}
}

func (d DepartmentHandelr) validateDelete(id int) int {
	res, err := d.dataStore.Delete(id)
	if err != nil {
		return res
	}
	return res
}
