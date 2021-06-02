package handler

import (
	"bytes"
	errors2 "errors"
	"fmt"
	"net/http/httptest"
	"reflect"
	"testing"

	"developer.zopsmart.com/go/gofr/examples/using-postgres/model"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
	"developer.zopsmart.com/go/gofr/pkg/gofr/responder"
)

type mockStore struct{}

func (m mockStore) Get(c *gofr.Context) (*[]model.Customer, error) {
	p := c.Param("mock")
	if p == "success" {
		return nil, nil
	}

	return nil, errors2.New("error fetching customer listing")
}

func (m mockStore) GetByID(c *gofr.Context, id int) (*model.Customer, error) {
	if id == 1 {
		return &model.Customer{
			ID:   1,
			Name: "some name",
		}, nil
	}

	return nil, errors.EntityNotFound{
		Entity: "customer",
		ID:     fmt.Sprint(id),
	}
}

func (m mockStore) Update(c *gofr.Context, customer model.Customer) (*model.Customer, error) {
	if customer.Name == "some name" {
		return nil, nil
	}

	return nil, errors2.New("error updating customer")
}

func (m mockStore) Create(c *gofr.Context, customer model.Customer) (*model.Customer, error) {
	switch customer.Name {
	case "some name":
		return &model.Customer{
			ID:   1,
			Name: "success",
		}, nil
	case "mock body error":
		return nil, errors.InvalidParam{Param: []string{"body"}}
	case `{"id":1}`:
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	return nil, errors2.New("error adding new customer")
}

func (m mockStore) Delete(c *gofr.Context, id int) error {
	if c.PathParam("id") == "123" {
		return nil
	}

	return errors2.New("error deleting customer")
}

func TestModel_AddCustomer(t *testing.T) {
	m := New(mockStore{})

	k := gofr.New()

	tests := []struct {
		body        []byte
		expectedErr error
	}{
		{
			body:        []byte(`{"id":1}`),
			expectedErr: errors.InvalidParam{Param: []string{"id"}},
		},
		{
			body:        []byte(`{"name":"some name"}`),
			expectedErr: nil,
		},
		{
			body:        []byte(`mock body error`),
			expectedErr: errors.InvalidParam{Param: []string{"body"}},
		},
		{
			body:        []byte(`{"name":"creation error"}`),
			expectedErr: errors2.New("error adding new customer"),
		},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "http://dummy", bytes.NewReader(test.body))

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)
		c := gofr.NewContext(res, req, k)

		_, gotErr := m.Create(c)
		if !reflect.DeepEqual(gotErr, test.expectedErr) {
			t.Errorf("FAILED, Expected: %v, Got: %v", test.expectedErr, gotErr)
		}
	}
}

func TestModel_UpdateCustomer(t *testing.T) {
	m := New(mockStore{})

	k := gofr.New()

	tests := []struct {
		body         []byte
		expectedErr  error
		urlPathParam string
	}{
		{
			body:        nil,
			expectedErr: errors.MissingParam{Param: []string{"id"}},
		},
		{
			body:         nil,
			urlPathParam: "abc123",
			expectedErr:  errors.InvalidParam{Param: []string{"id"}},
		},
		{
			body:         []byte(`mock body error`),
			expectedErr:  errors.InvalidParam{Param: []string{"body"}},
			urlPathParam: "123",
		},
		{
			body:         []byte(`{"name":"some name"}`),
			urlPathParam: "123",
		},
		{
			body:         []byte(`{"name":"creation error"}`),
			expectedErr:  errors2.New("error updating customer"),
			urlPathParam: "123",
		},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "http://dummy", bytes.NewReader(test.body))

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		c := gofr.NewContext(res, req, k)

		c.SetPathParams(map[string]string{
			"id": test.urlPathParam,
		})

		_, gotErr := m.Update(c)
		if !reflect.DeepEqual(gotErr, test.expectedErr) {
			t.Errorf("FAILED, Expected: %v, Got: %v", test.expectedErr, gotErr)
		}
	}
}

func TestModel_GetCustomerById(t *testing.T) {
	m := New(mockStore{})

	k := gofr.New()

	tests := []struct {
		id           int
		expectedErr  error
		expectedResp *model.Customer
		urlPathParam string
	}{
		{
			id: 1,
			expectedResp: &model.Customer{
				ID:   1,
				Name: "some name",
			},
			urlPathParam: "1",
		},
		{
			id:           1,
			urlPathParam: "absd123",
			expectedErr:  errors.InvalidParam{Param: []string{"id"}},
		},
		{
			id:          1,
			expectedErr: errors.MissingParam{Param: []string{"id"}},
		},
		{
			id: 2,
			expectedErr: errors.EntityNotFound{
				Entity: "customer",
				ID:     "2",
			},
			urlPathParam: "2",
		},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "http://dummy", nil)

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		c := gofr.NewContext(res, req, k)

		c.SetPathParams(map[string]string{
			"id": test.urlPathParam,
		})

		_, gotErr := m.GetByID(c)
		if !reflect.DeepEqual(gotErr, test.expectedErr) {
			t.Errorf("FAILED, Expected: %v, Got: %v", test.expectedErr, gotErr)
		}
	}
}

func TestModel_DeleteCustomer(t *testing.T) {
	m := New(mockStore{})

	k := gofr.New()

	tests := []struct {
		id           int
		expectedErr  error
		urlPathParam string
	}{
		{
			id:           123,
			urlPathParam: "123",
		},
		{
			id:           1234,
			urlPathParam: "1234",
			expectedErr:  errors2.New("error deleting customer"),
		},
		{
			id:           1,
			urlPathParam: "absd123",
			expectedErr:  errors.InvalidParam{Param: []string{"id"}},
		},
		{
			id:          1,
			expectedErr: errors.MissingParam{Param: []string{"id"}},
		},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "http://dummy", nil)

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		c := gofr.NewContext(res, req, k)

		c.SetPathParams(map[string]string{
			"id": test.urlPathParam,
		})

		_, gotErr := m.Delete(c)
		if !reflect.DeepEqual(gotErr, test.expectedErr) {
			t.Errorf("FAILED, Expected: %v, Got: %v", test.expectedErr, gotErr)
		}
	}
}

func TestModel_GetCustomers(t *testing.T) {
	m := New(mockStore{})

	k := gofr.New()

	tests := []struct {
		expectedErr  error
		mockParamStr string
	}{
		{
			mockParamStr: "mock=success",
		},
		{
			expectedErr: errors2.New("error fetching customer listing"),
		},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "http://dummy?"+test.mockParamStr, nil)

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)
		c := gofr.NewContext(res, req, k)

		_, gotErr := m.Get(c)
		if !reflect.DeepEqual(gotErr, test.expectedErr) {
			t.Errorf("FAILED, Expected: %v, Got: %v", test.expectedErr, gotErr)
		}
	}
}
