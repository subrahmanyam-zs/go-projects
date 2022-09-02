package department

import (
	"strings"

	"golang.org/x/exp/slices"

	"developer.zopsmart.com/go/gofr/Emp-Dept/datastore"
	"developer.zopsmart.com/go/gofr/Emp-Dept/entities"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Service struct {
	store datastore.Department
}

func New(store datastore.Department) Service {
	return Service{store: store}
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

func validation(dept entities.Department) errors.InvalidParam {
	switch {
	case !validateName(dept.DeptName):
		return errors.InvalidParam{Param: []string{"name"}}
	case !validateFloorNo(dept.FloorNo):
		return errors.InvalidParam{Param: []string{"floor no"}}
	}

	return errors.InvalidParam{}
}

func (s Service) Post(ctx *gofr.Context, dept entities.Department) (interface{}, error) {
	if err := validation(dept); err.Param != nil {
		return nil, err
	}

	res, err := s.store.Post(ctx, dept)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s Service) Put(ctx *gofr.Context, id int, dept entities.Department) (interface{}, error) {
	if err := validation(dept); err.Param != nil {
		return nil, err
	}

	res, err := s.store.Put(ctx, id, dept)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s Service) Delete(ctx *gofr.Context, id int) (int, error) {
	res, err := s.store.Delete(ctx, id)
	if err != nil {
		return res, err
	}

	return res, nil
}
