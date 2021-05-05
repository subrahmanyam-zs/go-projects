package handlers

import (
	"strconv"

	"github.com/zopsmart/gofr/examples/using-ycql/entity"
	"github.com/zopsmart/gofr/examples/using-ycql/store"
	"github.com/zopsmart/gofr/pkg/errors"
	"github.com/zopsmart/gofr/pkg/gofr"
)

type Shop struct {
	model store.Shop
}

func New(ps store.Shop) Shop {
	return Shop{
		model: ps,
	}
}
func (s Shop) Get(c *gofr.Context) (interface{}, error) {
	var filter entity.Shop

	params := c.Params()

	filter.ID, _ = strconv.Atoi(params["id"])
	filter.Name = params["name"]
	filter.Location = params["location"]
	filter.State = params["state"]

	return s.model.Get(c, filter), nil
}

func (s Shop) Create(c *gofr.Context) (interface{}, error) {
	var data entity.Shop
	// json error
	if err := c.Bind(&data); err != nil {
		return nil, err
	}

	records := s.model.Get(c, entity.Shop{ID: data.ID})

	if len(records) > 0 {
		return nil, errors.EntityAlreadyExists{}
	}

	return s.model.Create(c, data)
}

func (s Shop) Delete(c *gofr.Context) (interface{}, error) {
	var filter entity.Shop

	filter.ID, _ = strconv.Atoi(c.PathParam("id"))

	id := c.PathParam("id")
	// first verify that value exit or not?
	records := s.model.Get(c, filter)

	if len(records) == 0 {
		return nil, errors.EntityNotFound{Entity: "person", ID: id}
	}

	return nil, s.model.Delete(c, c.PathParam("id"))
}

func (s Shop) Update(c *gofr.Context) (interface{}, error) {
	var data entity.Shop

	if err := c.Bind(&data); err != nil {
		return nil, err
	}

	data.ID, _ = strconv.Atoi(c.PathParam("id"))
	records := s.model.Get(c, entity.Shop{ID: data.ID})

	if len(records) == 0 {
		return nil, errors.EntityNotFound{Entity: "person", ID: c.PathParam("id")}
	}

	return s.model.Update(c, data)
}
