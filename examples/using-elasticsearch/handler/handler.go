package handler

import (
	"net/http"

	"github.com/zopsmart/gofr/examples/using-elasticsearch/model"
	"github.com/zopsmart/gofr/examples/using-elasticsearch/store"
	"github.com/zopsmart/gofr/pkg/errors"
	"github.com/zopsmart/gofr/pkg/gofr"
)

type Customer struct {
	store store.Customer
}

func New(c store.Customer) *Customer {
	return &Customer{store: c}
}

func (c Customer) Index(context *gofr.Context) (interface{}, error) {
	name := context.Param("name")

	resp, err := c.store.Get(context, name)
	if err != nil {
		return nil, &errors.Response{StatusCode: http.StatusInternalServerError, Reason: "something unexpected happened"}
	}

	return resp, nil
}

func (c Customer) Read(context *gofr.Context) (interface{}, error) {
	id := context.PathParam("id")

	if id == "" {
		return nil, errors.MissingParam{Param: []string{"id"}}
	}

	resp, err := c.store.GetByID(context, id)
	if err != nil {
		return nil, &errors.Response{StatusCode: http.StatusInternalServerError, Reason: "something unexpected happened"}
	}

	return resp, nil
}

func (c Customer) Update(context *gofr.Context) (interface{}, error) {
	id := context.PathParam("id")
	if id == "" {
		return nil, errors.MissingParam{Param: []string{"id"}}
	}

	var cust model.Customer

	if err := context.Bind(&cust); err != nil {
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}

	resp, err := c.store.Update(context, cust, id)
	if err != nil {
		return nil, &errors.Response{StatusCode: http.StatusInternalServerError, Reason: "something unexpected happened"}
	}

	return resp, nil
}

func (c Customer) Create(context *gofr.Context) (interface{}, error) {
	var cust model.Customer
	if err := context.Bind(&cust); err != nil {
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}

	resp, err := c.store.Create(context, cust)
	if err != nil {
		return nil, &errors.Response{StatusCode: http.StatusInternalServerError, Reason: "something unexpected happened"}
	}

	return resp, nil
}
func (c Customer) Delete(context *gofr.Context) (interface{}, error) {
	i := context.PathParam("id")
	if i == "" {
		return nil, errors.MissingParam{Param: []string{"id"}}
	}

	if err := c.store.Delete(context, i); err != nil {
		return nil, &errors.Response{StatusCode: http.StatusInternalServerError, Reason: "something unexpected happened"}
	}

	return "Deleted successfully", nil
}
