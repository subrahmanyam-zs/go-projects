package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"developer.zopsmart.com/go/gofr/examples/data-layer-with-postgres/model"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
	"developer.zopsmart.com/go/gofr/pkg/gofr/responder"

	"github.com/stretchr/testify/assert"
)

type mockStore struct{}

func (m mockStore) Get(ctx *gofr.Context) ([]model.Customer, error) {
	p := ctx.Param("mock")
	if p == "success" {
		return nil, nil
	}

	return nil, errors.Error("error fetching customer listing")
}

func (m mockStore) GetByID(ctx *gofr.Context, id int) (model.Customer, error) {
	if id == 1 {
		return model.Customer{ID: 1, Name: "some name"}, nil
	}

	return model.Customer{}, errors.EntityNotFound{Entity: "customer", ID: fmt.Sprint(id)}
}

func (m mockStore) Update(ctx *gofr.Context, customer model.Customer) (model.Customer, error) {
	if customer.Name == "some name" {
		return model.Customer{}, nil
	}

	return model.Customer{}, errors.Error("error updating customer")
}

func (m mockStore) Create(ctx *gofr.Context, customer model.Customer) (model.Customer, error) {
	switch customer.Name {
	case "some name":
		return model.Customer{ID: 1, Name: "success"}, nil
	case "mock body error":
		return model.Customer{}, errors.InvalidParam{Param: []string{"body"}}
	case `{"id":1}`:
		return model.Customer{}, errors.InvalidParam{Param: []string{"id"}}
	}

	return model.Customer{}, errors.Error("error adding new customer")
}

func (m mockStore) Delete(ctx *gofr.Context, id int) error {
	if ctx.PathParam("id") == "123" {
		return nil
	}

	return errors.Error("error deleting customer")
}

func TestModel_AddCustomer(t *testing.T) {
	h := New(mockStore{})

	app := gofr.New()

	tests := []struct {
		desc string
		body []byte
		err  error
	}{
		{"create with invalid id", []byte(`{"id":1}`), errors.InvalidParam{Param: []string{"id"}}},
		{"create succuss", []byte(`{"name":"some name"}`), nil},
		{"create invalid body", []byte(`mock body error`), errors.InvalidParam{Param: []string{"body"}}},
		{"create error", []byte(`{"name":"creation error"}`), errors.Error("error adding new customer")},
	}

	for i, tc := range tests {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "http://dummy", bytes.NewReader(tc.body))

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)
		ctx := gofr.NewContext(res, req, app)

		_, err := h.Create(ctx)
		assert.Equal(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

func TestModel_UpdateCustomer(t *testing.T) {
	h := New(mockStore{})

	app := gofr.New()

	tests := []struct {
		desc string
		body []byte
		err  error
		id   string
	}{
		{"missing id", nil, errors.MissingParam{Param: []string{"id"}}, ""},
		{"invalid id", nil, errors.InvalidParam{Param: []string{"id"}}, "abc123"},
		{"invalid body", []byte(`{`), errors.InvalidParam{Param: []string{"body"}}, "123"},
		{"update succuss", []byte(`{"name":"some name"}`), nil, "123"},
		{"update error", []byte(`{"name":"creation error"}`), errors.Error("error updating customer"), "123"},
	}

	for i, tc := range tests {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "http://dummy", bytes.NewReader(tc.body))

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(res, req, app)

		ctx.SetPathParams(map[string]string{
			"id": tc.id,
		})

		_, err := h.Update(ctx)
		assert.Equal(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

func TestModel_GetCustomerById(t *testing.T) {
	h := New(mockStore{})

	app := gofr.New()

	tests := []struct {
		desc string
		id   string
		resp interface{}
		err  error
	}{
		{"get by id succuss", "1", model.Customer{ID: 1, Name: "some name"}, nil},
		{"invalid id", "absd123", nil, errors.InvalidParam{Param: []string{"id"}}},
		{"missing id", "", nil, errors.MissingParam{Param: []string{"id"}}},
		{"id not found", "2", nil, errors.EntityNotFound{Entity: "customer", ID: "2"}},
	}

	for i, tc := range tests {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "http://dummy", nil)

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(res, req, app)

		ctx.SetPathParams(map[string]string{
			"id": tc.id,
		})

		resp, err := h.GetByID(ctx)
		assert.Equal(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)

		assert.Equal(t, tc.resp, resp, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

func TestModel_DeleteCustomer(t *testing.T) {
	h := New(mockStore{})

	app := gofr.New()

	tests := []struct {
		desc string
		id   string
		err  error
	}{
		{"delete succuss", "123", nil},
		{"delete fail", "1234", errors.Error("error deleting customer")},
		{"invalid id", "absd123", errors.InvalidParam{Param: []string{"id"}}},
		{"missing id", "", errors.MissingParam{Param: []string{"id"}}},
	}

	for i, tc := range tests {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "http://dummy", nil)

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(res, req, app)

		ctx.SetPathParams(map[string]string{
			"id": tc.id,
		})

		_, err := h.Delete(ctx)
		assert.Equal(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

func TestModel_GetCustomers(t *testing.T) {
	h := New(mockStore{})

	app := gofr.New()

	tests := []struct {
		desc         string
		mockParamStr string
		err          error
	}{
		{"get success", "mock=success", nil},
		{"get fail", "", errors.Error("error fetching customer listing")},
	}

	for i, tc := range tests {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "http://dummy?"+tc.mockParamStr, nil)

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)
		ctx := gofr.NewContext(res, req, app)

		_, err := h.Get(ctx)
		assert.Equal(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}
