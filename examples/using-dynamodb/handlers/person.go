package handlers

import (
	"developer.zopsmart.com/go/gofr/examples/using-dynamodb/model"
	"developer.zopsmart.com/go/gofr/examples/using-dynamodb/store"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type handler struct {
	store store.Person
}

// nolint:gocritic //exporting return value is not necessary here
// New factory function for person handler
func New(store store.Person) handler {
	return handler{store: store}
}

func (h handler) Create(c *gofr.Context) (interface{}, error) {
	var person model.Person

	err := c.Bind(&person)
	if err != nil {
		return nil, err
	}

	err = h.store.Create(c, person)
	if err != nil {
		return nil, err
	}

	return "Successful", nil
}

func (h handler) GetByID(c *gofr.Context) (interface{}, error) {
	id := c.PathParam("id")

	person, err := h.store.Get(c, id)
	if err != nil {
		return nil, err
	}

	return person, nil
}

func (h handler) Update(c *gofr.Context) (interface{}, error) {
	id := c.PathParam("id")

	var person model.Person

	err := c.Bind(&person)
	if err != nil {
		return nil, err
	}

	person.ID = id

	err = h.store.Update(c, person)
	if err != nil {
		return nil, err
	}

	return "Successful", nil
}

func (h handler) Delete(c *gofr.Context) (interface{}, error) {
	id := c.PathParam("id")

	return nil, h.store.Delete(c, id)
}
