package handler

import (
	"developer.zopsmart.com/go/gofr/examples/mock-c-layer/models"
	"developer.zopsmart.com/go/gofr/examples/mock-c-layer/store"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type handler struct {
	store store.Brand
}

// New is factory function for handler layer
//nolint:revive // handler should not be used without proper initilization with required dependency
func New(b store.Brand) handler {
	return handler{store: b}
}

// Get interacts with core layer to get brands
func (h handler) Get(ctx *gofr.Context) (interface{}, error) {
	resp, err := h.store.Get(ctx)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Create interacts with core layer to create a Model
func (h handler) Create(ctx *gofr.Context) (interface{}, error) {
	var b models.Brand

	err := ctx.Bind(&b)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"request body"}}
	}

	resp, err := h.store.Create(ctx, b)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
