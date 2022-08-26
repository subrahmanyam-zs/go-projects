package employee

import (
	"context"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"

	entities "EmployeeDepartment/entities"
	"EmployeeDepartment/errorsHandler"
	"EmployeeDepartment/store"
)

type Handler struct {
	dataStore store.Employee
}

func New(emp store.Employee) Handler {
	return Handler{dataStore: emp}
}

func validateDob(dob string) bool {
	re := regexp.MustCompile(`(0?[1-9]|[12]\d|3[01])-(0?[1-9]|1[012])-((19|20)\d\d)`)
	if !re.MatchString(dob) {
		return false
	}

	sliceYr := strings.Split(dob, "-")

	year, err := strconv.Atoi(sliceYr[2])
	if err != nil {
		return false
	}

	const yr = 1999
	if (re.MatchString(dob)) && (year < yr) {
		return true
	}

	return false
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

func validateDepartment(majors string, department entities.Department) bool {
	department.Name = strings.ToUpper(department.Name)

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

	default:
		if department.Name == "HR" {
			return true
		}

		return false
	}
}

func validation(emp *entities.Employee, dept entities.Department) errorsHandler.InvalidDetails {
	switch {
	case !validateDob(emp.Dob):
		return errorsHandler.InvalidDetails{Msg: "Dob age should be greater than 22"}
	case !validateMajors(emp.Majors):
		return errorsHandler.InvalidDetails{Msg: "Major"}
	case !validateCity(emp.City):
		return errorsHandler.InvalidDetails{Msg: "City"}
	case !validateDepartment(emp.Majors, dept):
		return errorsHandler.InvalidDetails{Msg: "Did"}
	}

	return errorsHandler.InvalidDetails{}
}

func (h Handler) Create(ctx context.Context, employee *entities.Employee) (*entities.Employee, error) {
	dept, err := h.GetDepartment(ctx, employee.DId)
	if err != nil {
		return &entities.Employee{}, err
	}
	if err := validation(employee, dept); err.Msg != "" {
		return &entities.Employee{}, &err
	}

	res, err := h.dataStore.Create(ctx, employee)
	if err != nil {
		return &entities.Employee{}, err
	}

	return res, nil
}

func (h Handler) Update(ctx context.Context, id uuid.UUID, employee *entities.Employee) (*entities.Employee, error) {
	dept, err := h.GetDepartment(ctx, employee.DId)
	if err != nil {
		return nil, err
	}
	if err := validation(employee, dept); err.Msg != "" {
		return &entities.Employee{}, &err
	}

	res, err := h.dataStore.Update(ctx, id, employee)
	if err != nil {
		return &entities.Employee{}, err
	}

	return res, nil
}

func (h Handler) Delete(ctx context.Context, id uuid.UUID) (int, error) {
	res, err := h.dataStore.Delete(ctx, id)
	if err != nil {
		return http.StatusNotFound, err
	}

	return res, nil
}

func (h Handler) Read(ctx context.Context, id uuid.UUID) (entities.EmployeeAndDepartment, error) {
	res, err := h.dataStore.Read(ctx, id)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (h Handler) ReadAll(para store.Parameters) ([]entities.EmployeeAndDepartment, error) {
	res, err := h.dataStore.ReadAll(para)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (h Handler) GetDepartment(ctx context.Context, id int) (entities.Department, error) {
	res, err := h.dataStore.ReadDepartment(ctx, id)
	if err != nil {
		return entities.Department{}, err
	}

	return res, nil
}
