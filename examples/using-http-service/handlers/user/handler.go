package user

import (
	"strings"

	"developer.zopsmart.com/go/gofr/examples/using-http-service/services"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type handler struct {
	service services.User
}

// New is factory function for handler layer
//nolint:revive // handler should not be used without proper initilization with required dependency
func New(service services.User) handler {
	return handler{service: service}
}

// Get is a handler function of type gofr.Handler that uses HTTP Service to make downstream calls
func (h handler) Get(ctx *gofr.Context) (interface{}, error) {
	name := ctx.PathParam("name")
	if strings.TrimSpace(name) == "" {
		return nil, errors.MissingParam{Param: []string{"name"}}
	}

	resp, err := h.service.Get(ctx, name)
	if err != nil {
		return nil, err // avoiding partial content response
	}

	return resp, nil
}
