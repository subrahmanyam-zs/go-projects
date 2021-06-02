package handler

import (
	"developer.zopsmart.com/go/gofr/examples/using-solr/store"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Customer struct {
	store store.Customer
}

// New initializes the consumer layer
func New(s store.Customer) *Customer {
	return &Customer{store: s}
}

// List lists the customers based on the parameters sent in the query
func (customer Customer) List(c *gofr.Context) (interface{}, error) {
	id := c.Param("id")
	if id == "" {
		return nil, errors.MissingParam{Param: []string{"id"}}
	}

	filter := store.Filter{ID: id, Name: c.Param("name")}

	resp, err := customer.store.List(c, "customer", filter)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Create creates a document in the customer collection
func (customer Customer) Create(c *gofr.Context) (interface{}, error) {
	var model store.Model

	err := c.Bind(&model)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}

	if model.Name == "" {
		return nil, errors.InvalidParam{Param: []string{"name"}}
	}

	return nil, customer.store.Create(c, "customer", model)
}

// Update updates a document in the customer collection
func (customer Customer) Update(c *gofr.Context) (interface{}, error) {
	var model store.Model

	err := c.Bind(&model)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}

	if model.Name == "" {
		return nil, errors.InvalidParam{Param: []string{"name"}}
	}

	if model.ID == 0 {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	return nil, customer.store.Update(c, "customer", model)
}

// Delete deletes a document in the customer collection
func (customer Customer) Delete(c *gofr.Context) (interface{}, error) {
	var model store.Model

	err := c.Bind(&model)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}

	if model.ID == 0 {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	return nil, customer.store.Delete(c, "customer", model)
}
