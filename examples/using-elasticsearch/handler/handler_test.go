package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
	"developer.zopsmart.com/go/gofr/examples/using-elasticsearch/model"
	errors2 "developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
)

type mockStore struct{}

func (m mockStore) Get(context *gofr.Context, name string) ([]model.Customer, error) {
	if name == "error" {
		return nil, errors.New("elastic search error")
	} else if name == "multiple" {
		return []model.Customer{{ID: "12", Name: "ee", City: "city"}, {ID: "189"}}, nil
	}

	return nil, nil
}

func (m mockStore) GetByID(context *gofr.Context, id string) (*model.Customer, error) {
	if id == "o978" {
		return nil, errors.New("error")
	}

	return &model.Customer{ID: "ipo897", Name: "Marc"}, nil
}

func (m mockStore) Update(context *gofr.Context, customer model.Customer, id string) (*model.Customer, error) {
	if id == "ofjru3343" {
		return nil, errors.New("error")
	}

	return &model.Customer{ID: "ipo897", Name: "Henry"}, nil
}

func (m mockStore) Create(context *gofr.Context, customer model.Customer) (*model.Customer, error) {
	if customer.Name == "March" {
		return nil, errors.New("cannot insert")
	}

	return &model.Customer{ID: "weop24444", Name: "Mike"}, nil
}

func (m mockStore) Delete(context *gofr.Context, id string) error {
	if id == "ef444" {
		return errors.New("error while deleting")
	}

	return nil
}

func TestCustomer_Index(t *testing.T) {
	testcases := []struct {
		name      string
		customers interface{}
		err       error
	}{
		{"error", nil, &errors2.Response{StatusCode: 500, Reason: "something unexpected happened"}},
		{"multiple", []model.Customer{{ID: "12", Name: "ee", City: "city"}, {ID: "189"}}, nil},
	}

	customer := New(mockStore{})

	for i, v := range testcases {
		k := gofr.New()
		req := httptest.NewRequest("GET", "/customer?name="+v.name, nil)
		r := request.NewHTTPRequest(req)
		context := gofr.NewContext(nil, r, k)
		resp, err := customer.Index(context)

		if !reflect.DeepEqual(v.err, err) {
			t.Errorf("[TESTCASE%d]Failed. Got: %v\tExpected: %v\n", i+1, err, v.err)
		}

		if !reflect.DeepEqual(resp, v.customers) {
			t.Errorf("[TESTCASE%d]Failed. Got: %v\tExpected: %v\n", i+1, resp, v.customers)
		}
	}
}

// nolint:dupl // some statement are similar to statement in other function
func TestCustomer_Read(t *testing.T) {
	testcases := []struct {
		id       string
		customer interface{}
		err      error
	}{
		{"", nil, errors2.MissingParam{Param: []string{"id"}}},
		{"o978", nil, &errors2.Response{StatusCode: 500, Reason: "something unexpected happened"}},
		{"ity6", &model.Customer{ID: "ipo897", Name: "Marc"}, nil},
	}

	customer := New(mockStore{})

	for i, v := range testcases {
		k := gofr.New()
		req := httptest.NewRequest("GET", "/customer?"+v.id, nil)
		req = mux.SetURLVars(req, map[string]string{
			"id": v.id,
		})

		r := request.NewHTTPRequest(req)
		context := gofr.NewContext(nil, r, k)
		resp, err := customer.Read(context)

		if !reflect.DeepEqual(err, v.err) {
			t.Errorf("[TESTCASE%d]Failed. Got: %v\tExpected: %v\n", i+1, err, v.err)
		}

		if !reflect.DeepEqual(resp, v.customer) {
			t.Errorf("[TESTCASE%d]Failed. Got: %v\tExpected: %v\n", i+1, resp, v.customer)
		}
	}
}

func TestCustomer_Create(t *testing.T) {
	testcases := []struct {
		customer model.Customer
		err      error
		resp     interface{}
	}{
		{model.Customer{Name: "March"}, &errors2.Response{StatusCode: 500, Reason: "something unexpected happened"}, nil},
		{model.Customer{Name: "Henry", City: "Marc City"}, nil, &model.Customer{ID: "weop24444", Name: "Mike"}},
	}

	customer := New(mockStore{})

	for i, v := range testcases {
		k := gofr.New()
		body, _ := json.Marshal(v.customer)
		req := httptest.NewRequest("GET", "/customer", bytes.NewBuffer(body))

		r := request.NewHTTPRequest(req)
		context := gofr.NewContext(nil, r, k)
		resp, err := customer.Create(context)

		if !reflect.DeepEqual(err, v.err) {
			t.Errorf("[TESTCASE%d]Failed. Got: %v\tExpected: %v\n", i+1, err, v.err)
		}

		if !reflect.DeepEqual(resp, v.resp) {
			t.Errorf("[TESTCASE%d]Failed. Got: %v\tExpected: %v\n", i+1, resp, v.resp)
		}
	}
}

func TestCustomer_Update(t *testing.T) {
	testcases := []struct {
		id       string
		customer interface{}
		err      error
		resp     *model.Customer
	}{
		{"", nil, errors2.MissingParam{Param: []string{"id"}}, nil},
		{"ofjru3343", nil, &errors2.Response{StatusCode: 500, Reason: "something unexpected happened"}, nil},
		{"ity6", &model.Customer{ID: "ipo897", Name: "Henry"}, nil, &model.Customer{ID: "ipo897", Name: "Marc"}},
	}

	customer := New(mockStore{})

	for i, v := range testcases {
		k := gofr.New()
		body, _ := json.Marshal(v.customer)
		req := httptest.NewRequest("GET", "/customer?"+v.id, bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{
			"id": v.id,
		})

		r := request.NewHTTPRequest(req)
		context := gofr.NewContext(nil, r, k)
		resp, err := customer.Update(context)

		if !reflect.DeepEqual(err, v.err) {
			t.Errorf("[TESTCASE%d]Failed. Got: %v\tExpected: %v\n", i+1, err, v.err)
		}

		if !reflect.DeepEqual(resp, v.customer) {
			t.Errorf("[TESTCASE%d]Failed. Got: %v\tExpected: %v\n", i+1, resp, v.customer)
		}
	}
}

// nolint:dupl // some statement are similar to statement in other function
func TestCustomer_Delete(t *testing.T) {
	testcases := []struct {
		id       string
		customer interface{}
		err      error
	}{
		{"", nil, errors2.MissingParam{Param: []string{"id"}}},
		{"12", "Deleted successfully", nil},
		{"ef444", nil, &errors2.Response{StatusCode: 500, Reason: "something unexpected happened"}},
	}

	customer := New(mockStore{})

	for i, v := range testcases {
		k := gofr.New()
		req := httptest.NewRequest("DELETE", "/customer?"+v.id, nil)
		req = mux.SetURLVars(req, map[string]string{
			"id": v.id,
		})

		r := request.NewHTTPRequest(req)
		context := gofr.NewContext(nil, r, k)
		resp, err := customer.Delete(context)

		if !reflect.DeepEqual(err, v.err) {
			t.Errorf("[TESTCASE%d]Failed. Got: %v\tExpected: %v\n", i+1, err, v.err)
		}

		if !reflect.DeepEqual(resp, v.customer) {
			t.Errorf("[TESTCASE%d]Failed. Got: %v\tExpected: %v\n", i+1, resp, v.customer)
		}
	}
}

func TestInvalidBody(t *testing.T) {
	k := gofr.New()
	customer := New(mockStore{})
	req := httptest.NewRequest("POST", "/customer", nil)
	r := request.NewHTTPRequest(req)
	context := gofr.NewContext(nil, r, k)
	_, err := customer.Create(context)
	expectedErr := errors2.InvalidParam{Param: []string{"body"}}

	if !reflect.DeepEqual(err, expectedErr) {
		t.Errorf("failed for invalid body")
	}
}

func TestInvalidBodyUpdate(t *testing.T) {
	k := gofr.New()
	customer := New(mockStore{})
	req := httptest.NewRequest("POST", "/customer/1", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": "1",
	})
	r := request.NewHTTPRequest(req)
	context := gofr.NewContext(nil, r, k)
	_, err := customer.Update(context)
	expectedErr := errors2.InvalidParam{Param: []string{"body"}}

	if !reflect.DeepEqual(err, expectedErr) {
		t.Errorf("failed for invalid body")
	}
}
