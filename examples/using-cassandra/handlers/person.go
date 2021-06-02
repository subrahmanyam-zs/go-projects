package handlers

import (
	"strconv"

	"developer.zopsmart.com/go/gofr/examples/using-cassandra/entity"
	"developer.zopsmart.com/go/gofr/examples/using-cassandra/store"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Person struct {
	model store.Person
}

func New(ps store.Person) Person {
	return Person{
		model: ps,
	}
}

func (p Person) Get(c *gofr.Context) (interface{}, error) {
	var filter entity.Person

	val := c.Params()

	filter.ID, _ = strconv.Atoi(val["id"])
	filter.Name = val["name"]
	filter.Age, _ = strconv.Atoi(val["age"])
	filter.State = val["state"]

	return p.model.Get(c, filter), nil
}

func (p Person) Create(c *gofr.Context) (interface{}, error) {
	var data entity.Person
	// json error
	if err := c.Bind(&data); err != nil {
		return nil, err
	}

	records := p.model.Get(c, entity.Person{ID: data.ID})

	if len(records) > 0 {
		return nil, errors.EntityAlreadyExists{}
	}

	results, err := p.model.Create(c, data)

	return results, err
}

func (p Person) Delete(c *gofr.Context) (interface{}, error) {
	var filter entity.Person

	filter.ID, _ = strconv.Atoi(c.PathParam("id"))

	id := c.PathParam("id")
	// first verify that value exit or not?
	records := p.model.Get(c, filter)

	if len(records) == 0 {
		return nil, errors.EntityNotFound{Entity: "person", ID: id}
	}

	err := p.model.Delete(c, c.PathParam("id"))

	return nil, err
}

func (p Person) Update(c *gofr.Context) (interface{}, error) {
	var data entity.Person

	if err := c.Bind(&data); err != nil {
		return nil, err
	}

	data.ID, _ = strconv.Atoi(c.PathParam("id"))
	records := p.model.Get(c, entity.Person{ID: data.ID})

	if len(records) == 0 {
		return nil, errors.EntityNotFound{Entity: "person", ID: c.PathParam("id")}
	}

	return p.model.Update(c, data)
}
