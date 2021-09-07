package handlers

import (
	"developer.zopsmart.com/go/gofr/examples/using-dynamodb/model"
	"developer.zopsmart.com/go/gofr/examples/using-dynamodb/store"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type person struct {
	store store.Person
}

// New factory function for person handler
func New(store store.Person) person {
	return person{store: store}
}

func (p person) Create(c *gofr.Context) (interface{}, error) {
	var person model.Person

	err := c.Bind(&person)
	if err != nil {
		return nil, err
	}

	err = p.store.Create(c, person)
	if err != nil {
		return nil, err
	}

	return "Successful", nil
}

func (p person) GetByID(c *gofr.Context) (interface{}, error) {
	id := c.PathParam("id")

	person, err := p.store.Get(c, id)
	if err != nil {
		return nil, err
	}

	return person, nil
}

func (p person) Update(c *gofr.Context) (interface{}, error) {
	id := c.PathParam("id")

	var person model.Person

	err := c.Bind(&person)
	if err != nil {
		return nil, err
	}

	person.ID = id

	err = p.store.Update(c, person)
	if err != nil {
		return nil, err
	}

	return "Successful", nil
}

func (p person) Delete(c *gofr.Context) (interface{}, error) {
	id := c.PathParam("id")

	return nil, p.store.Delete(c, id)
}
