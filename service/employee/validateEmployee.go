package employee

import (
	entities2 "EmployeeDepartment/entities"
	"EmployeeDepartment/store"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type EmployeeHandler struct {
	datastore store.Employee
}

func New(emp store.Employee) EmployeeHandler {
	return EmployeeHandler{datastore: emp}
}

func validateId(id uuid.UUID) bool {
	if len(id.String()) != 36 {
		return false
	}
	return true
}

func validateDob(dob string) bool {
	re := regexp.MustCompile("(0?[1-9]|[12][0-9]|3[01])-(0?[1-9]|1[012])-((19|20)\\d\\d)")
	slice_yr := strings.Split(dob, "-")
	year, _ := strconv.Atoi(slice_yr[2])
	if (re.MatchString(dob)) && (year < 1999) {
		return true
	}
	return false
}

func validateCity(s string) bool {
	cities := []string{"Bangalore", "Kochi", "Mysore"}
	if slices.Contains(cities, s) {
		return true
	}
	return false
}

func validateMajors(s string) bool {
	s = strings.ToUpper(s)
	majors := []string{"CSE", "MCA", "MBA", "B.Com", "CA"}
	if slices.Contains(majors, s) {
		return true
	}
	return false
}
func validateDepartment(majors string, department entities2.Department) bool {
	fmt.Println(department)
	department.Name = strings.ToUpper(department.Name)
	fmt.Println(department.Name, majors)
	switch majors {
	case "CSE", "MCA":
		if department.Name == "TECH" {
			return true
		}
		return false
	case "B.COM", "CA":
		if department.Name == "ACCOUNTS" {
			return true
		}
		return false
	case "MBA":
		if department.Name == "HR" {
			return true
		}
		return false
	}
	return false
}

func (e EmployeeHandler) GetDepartment(id int) entities2.Department {
	res, err := e.datastore.ReadDepartment(id)
	if err != nil {
		return res
	}
	return res
}

func (e EmployeeHandler) Create(employee entities2.Employee) (entities2.Employee, error) {
	dept := e.GetDepartment(employee.DId)
	if validateDob(employee.Dob) && validateCity(employee.City) && validateMajors(employee.Majors) && validateDepartment(employee.Majors, dept) {
		res, err := e.datastore.Create(employee)
		if err != nil {
			return entities2.Employee{}, err
		}
		return res, nil
	}
	return entities2.Employee{}, errors.New("error")
}

func (e EmployeeHandler) Update(id uuid.UUID, employee entities2.Employee) (entities2.Employee, error) {
	fmt.Println(validateDob(employee.Dob), validateCity(employee.City), validateMajors(employee.Majors))
	if validateId(employee.Id) && validateDob(employee.Dob) && validateCity(employee.City) && validateMajors(employee.Majors) {
		res, err := e.datastore.Update(id, employee)
		if err != nil {
			return entities2.Employee{}, err
		}
		return res, nil
	}
	return entities2.Employee{}, errors.New("err")
}

func (e EmployeeHandler) Delete(id uuid.UUID) (int, error) {
	if validateId(id) {
		res, err := e.datastore.Delete(id)
		if err != nil {
			return http.StatusNotFound, err
		}
		return res, nil
	}
	return http.StatusNotFound, errors.New("err")
}

func (e EmployeeHandler) Read(id uuid.UUID) (entities2.EmployeeAndDepartment, error) {
	if validateId(id) {
		res, err := e.datastore.Read(id)
		if err != nil {
			return res, err
		}
		return res, nil
	}
	return entities2.EmployeeAndDepartment{}, errors.New("err")
}

func (e EmployeeHandler) ReadAll(name string, includeDepartment bool) ([]entities2.EmployeeAndDepartment, error) {
	res, err := e.datastore.ReadAll(name, includeDepartment)
	if err != nil {
		return res, err
	}
	return res, nil
}
