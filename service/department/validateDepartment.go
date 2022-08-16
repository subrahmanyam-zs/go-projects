package department

import (
	"EmployeeDepartment/entities"
	"EmployeeDepartment/store"
	"errors"
	"golang.org/x/exp/slices"
	"strings"
)

type DepartmentHandelr struct {
	dataStore store.Department
}

func New(dept store.Department) DepartmentHandelr {
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

func (d DepartmentHandelr) Create(department entities.Department) (entities.Department, error) {
	if validateName(department.Name) && validateFloorNo(department.FloorNo) {
		res, err := d.dataStore.Create(department)
		if err != nil {
			return entities.Department{}, err
		}
		return res, nil
	}
	return entities.Department{}, errors.New("error")
}

func (d DepartmentHandelr) Update(id int, department entities.Department) (entities.Department, error) {
	if validateName(department.Name) && validateFloorNo(department.FloorNo) {
		res, err := d.dataStore.Update(id, department)
		if err != nil {
			return res, err
		}
		return res, nil
	}
	return entities.Department{}, errors.New("error")
}

func (d DepartmentHandelr) Delete(id int) (int, error) {
	res, err := d.dataStore.Delete(id)
	if err != nil {
		return res, err
	}
	return res, nil
}
