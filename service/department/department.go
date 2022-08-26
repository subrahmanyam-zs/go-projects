package department

import (
	"context"
	"strings"

	"golang.org/x/exp/slices"

	"EmployeeDepartment/entities"
	"EmployeeDepartment/errorsHandler"
	"EmployeeDepartment/store"
)

type Handler struct {
	dataStore store.Department
}

func New(dept store.Department) Handler {
	return Handler{dataStore: dept}
}

func validateName(name string) bool {
	name = strings.ToUpper(name)
	deptNames := []string{"HR", "TECH", "ACCOUNTS"}

	return slices.Contains(deptNames, name)
}

func validateFloorNo(floorNo int) bool {
	if floorNo >= 1 && floorNo <= 3 {
		return true
	}

	return false
}

func validation(dept entities.Department) errorsHandler.InvalidDetails {
	switch {
	case !validateName(dept.Name):
		return errorsHandler.InvalidDetails{Msg: "name"}
	case !validateFloorNo(dept.FloorNo):
		return errorsHandler.InvalidDetails{Msg: "floorNo"}
	}

	return errorsHandler.InvalidDetails{}
}

func (h Handler) Create(ctx context.Context, department entities.Department) (entities.Department, error) {
	if err := validation(department); err.Msg != "" {
		return entities.Department{}, &err
	}

	res, err := h.dataStore.Create(ctx, department)
	if err != nil {
		return entities.Department{}, err
	}

	return res, nil
}

func (h Handler) Update(ctx context.Context, id int, department entities.Department) (entities.Department, error) {
	if err := validation(department); err.Msg != "" {
		return entities.Department{}, &err
	}

	res, err := h.dataStore.Update(ctx, id, department)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (h Handler) Delete(ctx context.Context, id int) (int, error) {
	res, err := h.dataStore.Delete(ctx, id)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (h Handler) GetDepartment(ctx context.Context, id int) (entities.Department, error) {
	dept, err := h.dataStore.GetDepartment(ctx, id)
	if err != nil {
		return dept, &errorsHandler.IDNotFound{Msg: "Id not found"}
	}
	return dept, nil
}
