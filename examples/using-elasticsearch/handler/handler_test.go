package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"developer.zopsmart.com/go/gofr/examples/using-elasticsearch/model"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
)

type mockStore struct{}

func (m mockStore) Get(_ *gofr.Context, name string) ([]model.Customer, error) {
	if name == "error" {
		return nil, errors.Error("elastic search error")
	} else if name == "multiple" {
		return []model.Customer{{ID: "12", Name: "ee", City: "city"}, {ID: "189"}}, nil
	}

	return nil, nil
}

func (m mockStore) GetByID(_ *gofr.Context, id string) (model.Customer, error) {
	if id == "o978" {
		return model.Customer{}, errors.Error("error")
	}

	return model.Customer{ID: "ipo897", Name: "Marc"}, nil
}

func (m mockStore) Update(_ *gofr.Context, _ model.Customer, id string) (model.Customer, error) {
	if id == "ofjru3343" {
		return model.Customer{}, errors.Error("error")
	}

	return model.Customer{ID: "ipo897", Name: "Henry"}, nil
}

func (m mockStore) Create(_ *gofr.Context, customer model.Customer) (model.Customer, error) {
	if customer.Name == "March" {
		return model.Customer{}, errors.Error("cannot insert")
	}

	return model.Customer{ID: "weop24444", Name: "Mike"}, nil
}

func (m mockStore) Delete(_ *gofr.Context, id string) error {
	if id == "ef444" {
		return errors.Error("error while deleting")
	}

	return nil
}

func TestCustomer_Index(t *testing.T) {
	testcases := []struct {
		name      string
		customers interface{}
		err       error
	}{
		{"error", nil, &errors.Response{StatusCode: 500, Reason: "something unexpected happened"}},
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

func TestCustomer_Read(t *testing.T) {
	testcases := []struct {
		id       string
		customer interface{}
		err      error
	}{
		{"", nil, errors.MissingParam{Param: []string{"id"}}},
		{"o978", nil, &errors.Response{StatusCode: 500, Reason: "something unexpected happened"}},
		{"ity6", model.Customer{ID: "ipo897", Name: "Marc"}, nil},
	}

	customer := New(mockStore{})

	for i, v := range testcases {
		k := gofr.New()

		req := httptest.NewRequest("GET", "/customer", nil)
		r := request.NewHTTPRequest(req)
		ctx := gofr.NewContext(nil, r, k)

		ctx.SetPathParams(map[string]string{
			"id": v.id,
		})

		resp, err := customer.Read(ctx)
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
		{model.Customer{Name: "March"}, &errors.Response{StatusCode: 500, Reason: "something unexpected happened"}, nil},
		{model.Customer{Name: "Henry", City: "Marc City"}, nil, model.Customer{ID: "weop24444", Name: "Mike"}},
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
		resp     interface{}
	}{
		{"", nil, errors.MissingParam{Param: []string{"id"}}, nil},
		{"ofjru3343", nil, &errors.Response{StatusCode: 500, Reason: "something unexpected happened"}, nil},
		{"ity6", &model.Customer{ID: "ipo897", Name: "Henry"}, nil, model.Customer{ID: "ipo897", Name: "Henry"}},
	}

	customer := New(mockStore{})

	for i, v := range testcases {
		k := gofr.New()

		body, _ := json.Marshal(v.customer)

		req := httptest.NewRequest(http.MethodGet, "/customer", bytes.NewBuffer(body))
		r := request.NewHTTPRequest(req)
		ctx := gofr.NewContext(nil, r, k)

		ctx.SetPathParams(map[string]string{
			"id": v.id,
		})

		resp, err := customer.Update(ctx)

		if !reflect.DeepEqual(err, v.err) {
			t.Errorf("[TESTCASE%d]Failed. Got: %v\tExpected: %v\n", i+1, err, v.err)
		}

		if !reflect.DeepEqual(resp, v.resp) {
			t.Errorf("[TESTCASE%d]Failed. Got: %v\tExpected: %v\n", i+1, resp, v.resp)
		}
	}
}

func TestCustomer_Delete(t *testing.T) {
	testcases := []struct {
		id       string
		customer interface{}
		err      error
	}{
		{"", nil, errors.MissingParam{Param: []string{"id"}}},
		{"12", "Deleted successfully", nil},
		{"ef444", nil, &errors.Response{StatusCode: 500, Reason: "something unexpected happened"}},
	}

	customer := New(mockStore{})

	for i, v := range testcases {
		k := gofr.New()

		req := httptest.NewRequest(http.MethodDelete, "/customer?"+v.id, nil)
		r := request.NewHTTPRequest(req)
		ctx := gofr.NewContext(nil, r, k)

		ctx.SetPathParams(map[string]string{
			"id": v.id,
		})

		resp, err := customer.Delete(ctx)
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
	req := httptest.NewRequest(http.MethodPost, "/customer", nil)
	r := request.NewHTTPRequest(req)
	context := gofr.NewContext(nil, r, k)
	_, err := customer.Create(context)
	expectedErr := errors.InvalidParam{Param: []string{"body"}}

	if !reflect.DeepEqual(err, expectedErr) {
		t.Errorf("failed for invalid body")
	}
}

func TestInvalidBodyUpdate(t *testing.T) {
	k := gofr.New()
	customer := New(mockStore{})
	req := httptest.NewRequest(http.MethodPost, "/customer", nil)
	r := request.NewHTTPRequest(req)
	ctx := gofr.NewContext(nil, r, k)

	ctx.SetPathParams(map[string]string{
		"id": "1",
	})

	expectedErr := errors.InvalidParam{Param: []string{"body"}}

	_, err := customer.Update(ctx)
	if !reflect.DeepEqual(err, expectedErr) {
		t.Errorf("failed for invalid body")
	}
}
