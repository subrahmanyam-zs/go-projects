package handler

import (
	"developer.zopsmart.com/go/gofr/examples/universal-example/gofr-services/service"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type handler struct {
	Service service.Service
}

//nolint:golint //handler should not get accessed, without initialization of service.Service
func New(svc service.Service) handler {
	return handler{Service: svc}
}

func (h handler) Log(c *gofr.Context) (interface{}, error) {
	svc := c.Param("service")
	if svc == "" {
		return nil, errors.MissingParam{Param: []string{"service"}}
	}

	return h.Service.Log(c, svc)
}

func (h handler) Hello(c *gofr.Context) (interface{}, error) {
	name := c.Param("name")

	return h.Service.Hello(c, name)
}
