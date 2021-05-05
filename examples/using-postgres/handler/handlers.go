package handler

import (
	"strconv"

	"github.com/zopsmart/gofr/examples/using-postgres/model"
	"github.com/zopsmart/gofr/examples/using-postgres/store"
	"github.com/zopsmart/gofr/pkg/errors"
	"github.com/zopsmart/gofr/pkg/gofr"
)

type Customer struct {
	store store.Store
}

func New(s store.Store) *Customer {
	return &Customer{
		store: s,
	}
}

type response struct {
	Customers *[]model.Customer
}

func (m Customer) Get(c *gofr.Context) (interface{}, error) {
	resp, err := m.store.Get(c)
	if err != nil {
		return nil, err
	}

	r := response{Customers: resp}

	return r, nil
}

func (m Customer) GetByID(c *gofr.Context) (interface{}, error) {
	i := c.PathParam("id")
	if i == "" {
		return nil, errors.MissingParam{Param: []string{"id"}}
	}

	id, err := strconv.Atoi(i)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	resp, err := m.store.GetByID(c, id)
	if err != nil {
		return nil, errors.EntityNotFound{
			Entity: "customer",
			ID:     i,
		}
	}

	return resp, nil
}

func (m Customer) Create(c *gofr.Context) (interface{}, error) {
	var cust model.Customer
	if err := c.Bind(&cust); err != nil {
		c.Logger.Errorf("error in binding: %v", err)
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}

	if cust.ID != 0 {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	resp, err := m.store.Create(c, cust)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m Customer) Update(c *gofr.Context) (interface{}, error) {
	i := c.PathParam("id")
	if i == "" {
		return nil, errors.MissingParam{Param: []string{"id"}}
	}

	id, err := strconv.Atoi(i)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	var cust model.Customer
	if err := c.Bind(&cust); err != nil {
		c.Logger.Errorf("error in binding: %v", err)
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}

	cust.ID = id

	resp, err := m.store.Update(c, cust)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m Customer) Delete(c *gofr.Context) (interface{}, error) {
	i := c.PathParam("id")
	if i == "" {
		return nil, errors.MissingParam{Param: []string{"id"}}
	}

	id, err := strconv.Atoi(i)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	if err := m.store.Delete(c, id); err != nil {
		return nil, err
	}

	return "Deleted successfully", nil
}
