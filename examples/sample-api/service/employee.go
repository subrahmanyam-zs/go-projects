package service

import (
	"developer.zopsmart.com/go/gofr/examples/sample-api/datastore"
	"developer.zopsmart.com/go/gofr/examples/sample-api/entity"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"golang.org/x/exp/slices"
	"strings"
)

type Service struct {
	dataStore datastore.Employee
}

func New(dataStore datastore.Employee) Service {
	return Service{dataStore: dataStore}
}

func validateName(name string) bool {
	if name == "" {
		return false
	}

	return true
}

func validateMajors(s string) bool {
	s = strings.ToUpper(s)
	majors := []string{"CSE", "MCA", "MBA", "B.Com", "CA"}

	return slices.Contains(majors, s)
}

func (s Service) Post(ctx *gofr.Context, emp entity.Employee) (interface{}, error) {
	if !validateName(emp.Name) {
		return entity.Employee{}, errors.InvalidParam{Param: []string{"name"}}
	} else if !validateMajors(emp.Majors) {
		return entity.Employee{}, errors.InvalidParam{Param: []string{"majors"}}
	}

	res, err := s.dataStore.Post(ctx, emp)
	if err != nil {
		return entity.Employee{}, errors.DB{Err: err}
	}

	return res, nil
}

func (s Service) Put(ctx *gofr.Context, id string, emp entity.Employee) (interface{}, error) {
	if !validateName(emp.Name) {
		return entity.Employee{}, errors.InvalidParam{Param: []string{"name"}}
	} else if !validateMajors(emp.Majors) {
		return entity.Employee{}, errors.InvalidParam{Param: []string{"majors"}}
	}

	res, err := s.dataStore.Put(ctx, id, emp)
	if err != nil {
		return nil, errors.DB{Err: err}
	}

	return res, nil
}

func (s Service) Delete(ctx *gofr.Context, id string) (interface{}, error) {
	res, err := s.dataStore.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s Service) Get(ctx *gofr.Context, id string) (interface{}, error) {
	res, err := s.dataStore.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s Service) GetAll(ctx *gofr.Context) (interface{}, error) {
	res, err := s.dataStore.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return res, err
}
