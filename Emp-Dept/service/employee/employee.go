package employee

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"

	"developer.zopsmart.com/go/gofr/Emp-Dept/datastore"
	"developer.zopsmart.com/go/gofr/Emp-Dept/entities"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Service struct {
	dataStore datastore.Employee
}

func New(d datastore.Employee) Service {
	return Service{dataStore: d}
}

func validateDob(dob string) bool {
	re := regexp.MustCompile(`(0?[1-9]|[12]\d|3[01])-(0?[1-9]|1[012])-((19|20)\d\d)`)
	if !re.MatchString(dob) {
		return false
	}

	sliceYr := strings.Split(dob, "-")

	year, _ := strconv.Atoi(sliceYr[2])

	const yr = 1999

	return (re.MatchString(dob)) && (year < yr)
}

func validateCity(s string) bool {
	cities := []string{"Bangalore", "Kochi", "Mysore"}
	return slices.Contains(cities, s)
}

func validateMajors(s string) bool {
	s = strings.ToUpper(s)
	majors := []string{"CSE", "MCA", "MBA", "B.Com", "CA"}

	return slices.Contains(majors, s)
}

func validateDepartment(majors string, department *entities.Department) bool {
	department.DeptName = strings.ToUpper(department.DeptName)

	switch majors {
	case "CSE", "MCA":
		if department.DeptName == "TECH" {
			return true
		}

		return false
	case "B.COM", "CA":
		if department.DeptName == "ACCOUNTS" {
			return true
		}

		return false

	default:
		if department.DeptName == "HR" {
			return true
		}

		return false
	}
}

func validateAll(emp *entities.Employee, dept *entities.Department) errors.InvalidParam {
	switch {
	case !validateDob(emp.Dob):
		return errors.InvalidParam{Param: []string{"age should be greater than 22"}}
	case !validateMajors(emp.Majors):
		return errors.InvalidParam{Param: []string{"majors"}}
	case !validateCity(emp.City):
		return errors.InvalidParam{Param: []string{"city"}}
	case !validateDepartment(emp.Majors, dept):
		return errors.InvalidParam{Param: []string{"deptID"}}
	}

	return errors.InvalidParam{}
}

func (s Service) Post(ctx *gofr.Context, emp entities.Employee) (interface{}, error) {
	dept, err := s.GetDepartment(ctx, emp.DeptID)
	if err != nil {
		return nil, err
	}

	if err := validateAll(&emp, &dept); err.Param != nil {
		return nil, err
	}

	res, err := s.dataStore.Post(ctx, emp)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s Service) Put(ctx *gofr.Context, id uuid.UUID, emp entities.Employee) (interface{}, error) {
	dept, err := s.GetDepartment(ctx, emp.DeptID)
	if err != nil {
		return nil, err
	}

	if err := validateAll(&emp, &dept); err.Param != nil {
		return nil, err
	}

	res, err := s.dataStore.Put(ctx, id, emp)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s Service) Delete(ctx *gofr.Context, id uuid.UUID) (int, error) {
	res, err := s.dataStore.Delete(ctx, id)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (s Service) Get(ctx *gofr.Context, id uuid.UUID) (interface{}, error) {
	res, err := s.dataStore.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s Service) GetAll(ctx *gofr.Context, name string, includeDept bool) (interface{}, error) {
	res, err := s.dataStore.GetAll(ctx, name, includeDept)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s Service) GetDepartment(ctx *gofr.Context, id int) (entities.Department, error) {
	res, err := s.dataStore.GetDepartment(ctx, id)
	if err != nil {
		return entities.Department{}, err
	}

	return res, nil
}
