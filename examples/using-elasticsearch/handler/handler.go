package handler

import (
	"net/http"

	"developer.zopsmart.com/go/gofr/examples/using-elasticsearch/model"
	"developer.zopsmart.com/go/gofr/examples/using-elasticsearch/store"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type customer struct {
	store store.Customer
}

func New(c store.Customer) customer {
	return customer{store: c}
}

func (c customer) Index(context *gofr.Context) (interface{}, error) {
	name := context.Param("name")

	resp, err := c.store.Get(context, name)
	if err != nil {
		return nil, &errors.Response{StatusCode: http.StatusInternalServerError, Reason: "something unexpected happened"}
	}

	return resp, nil
}

func (c customer) Read(context *gofr.Context) (interface{}, error) {
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

func (c customer) Update(context *gofr.Context) (interface{}, error) {
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

func (c customer) Create(context *gofr.Context) (interface{}, error) {
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
func (c customer) Delete(context *gofr.Context) (interface{}, error) {
	i := context.PathParam("id")
	if i == "" {
		return nil, errors.MissingParam{Param: []string{"id"}}
	}

	if err := c.store.Delete(context, i); err != nil {
		return nil, &errors.Response{StatusCode: http.StatusInternalServerError, Reason: "something unexpected happened"}
	}

	return "Deleted successfully", nil
}
