package handler

import (
	"developer.zopsmart.com/go/gofr/examples/mock-c-layer/store"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Brand struct {
	core store.Brand
}

// New returns a Model consumer
func New(b store.Brand) Brand {
	return Brand{b}
}

// Get interacts with core layer to get brands
func (b Brand) Get(c *gofr.Context) (interface{}, error) {
	resp, err := b.core.Get(c)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Create interacts with core layer to create a Model
func (b Brand) Create(c *gofr.Context) (interface{}, error) {
	m := store.Model{}

	err := c.Bind(&m)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"request body"}}
	}

	resp, err := b.core.Create(c, m)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
