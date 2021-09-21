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

//nolint:revive //exporting return value not necessary as we will be using New() outside the pkg.
func New(service services.User) handler {
	return handler{service: service}
}

// Get is a handler function of type gofr.Handler that uses HTTP Service to make downstream calls
func (h handler) Get(c *gofr.Context) (interface{}, error) {
	name := c.PathParam("name")
	if strings.TrimSpace(name) == "" {
		return nil, errors.MissingParam{Param: []string{"name"}}
	}

	resp, err := h.service.Get(c, name)
	if err != nil {
		return nil, err // avoiding partial content response
	}

	return resp, nil
}
