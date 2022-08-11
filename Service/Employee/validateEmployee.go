package Employee

import (
	"EmployeeDepartment/Handler/Entities"
	"EmployeeDepartment/Store"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type EmployeeHandler struct {
	datastore Store.Employee
}

func New(emp Store.Employee) EmployeeHandler {
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
	majors := []string{"CSE", "MCA", "MBA", "B.Com", "CA"}
	if slices.Contains(majors, s) {
		return true
	}
	return false
}

func (e EmployeeHandler) validatePost(employee Entities.Employee) Entities.Employee {
	if validateId(employee.Id) && validateDob(employee.Dob) && validateCity(employee.City) && validateMajors(employee.Majors) {
		res, err := e.datastore.Create(employee)
		if err != nil {
			return Entities.Employee{}
		}
		return res
	}
	return Entities.Employee{}
}

func (e EmployeeHandler) validatePut(id string, employee Entities.Employee) Entities.Employee {
	if validateId(employee.Id) && validateDob(employee.Dob) && validateCity(employee.City) && validateMajors(employee.Majors) {
		res, err := e.datastore.Update(id, employee)
		if err != nil {
			return Entities.Employee{}
		}
		return res
	}
	return Entities.Employee{}
}

func (e EmployeeHandler) validateDelete(id string) int {
	uid := uuid.MustParse(id)
	if validateId(uid) {
		res, err := e.datastore.Delete(uid)
		if err != nil {
			return http.StatusNotFound
		}
		return res
	}
	return http.StatusNotFound
}

func (e EmployeeHandler) validateGetById(id string) Entities.Employee {
	uid := uuid.MustParse(id)
	if validateId(uid) {
		res, err := e.datastore.Read(uid)
		if err != nil {
			return Entities.Employee{}
		}
		return res
	}
	return Entities.Employee{}
}

func (e EmployeeHandler) validateGetAll(name string, includeDepartment bool) []Entities.EmployeeAndDepartment {
	res, err := e.datastore.ReadAll(name, includeDepartment)
	if err != nil {
		return res
	}
	return res
}
