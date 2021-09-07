package handlers

import (
	"encoding/json"
	errors2 "errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"

	"developer.zopsmart.com/go/gofr/examples/using-mongo/entity"
	"developer.zopsmart.com/go/gofr/examples/using-mongo/store"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
)

func initializeHandlersTest(t *testing.T) (*store.MockCustomer, Customer, *gofr.Gofr) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	customerStore := store.NewMockCustomer(ctrl)
	customer := New(customerStore)
	k := gofr.New()

	return customerStore, customer, k
}

//nolint:govet //in table driven tests we don't add the key in the struct
func TestCustomer_Get(t *testing.T) {
	testCases := []struct {
		queryParams  string
		expectedResp interface{}
		mockErr      error
	}{
		{"", []*entity.Customer{{"Ponting", 24, "Sydney"}}, nil},
		{"name=Tim", []*entity.Customer{{"Tim", 35, "Munich"}}, nil},
		{"name=1", nil, errors.InvalidParam{Param: []string{"name"}}},
	}

	customerStore, customer, k := initializeHandlersTest(t)

	for index, tc := range testCases {
		req := httptest.NewRequest(http.MethodGet, "/customer?"+tc.queryParams, nil)
		r := request.NewHTTPRequest(req)
		context2 := gofr.NewContext(nil, r, k)

		customerStore.EXPECT().Get(gomock.Any(), gomock.Any()).Return(tc.expectedResp, tc.mockErr)

		resp, err := customer.Get(context2)
		if !reflect.DeepEqual(err, tc.mockErr) {
			t.Errorf("Testcase[%v] Failed\tExpected %v\tGot %v\n", index, tc.mockErr, err)
		}

		if !reflect.DeepEqual(resp, tc.expectedResp) {
			t.Errorf("Testcase[%v] Failed\tExpected %v\tGot %v\n", index, tc.expectedResp, resp)
		}
	}
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors2.New("test error")
}

func TestCustomer_Create_Invalid_Input_Error(t *testing.T) {
	expErr := errors2.New("test error")

	_, customer, k := initializeHandlersTest(t)
	req := httptest.NewRequest("GET", "/dummy", errReader(0))
	r := request.NewHTTPRequest(req)
	context2 := gofr.NewContext(nil, r, k)

	_, err := customer.Create(context2)
	if !reflect.DeepEqual(err, expErr) {
		t.Errorf("Testcase Failed\tExpected %v\tGot %v\n", expErr, err)
	}
}

func TestCustomer_Create_Invalid_JSON(t *testing.T) {
	input := `{"name":"Pirlo","age":"42","city":"Turin"}`
	expErr := &json.UnmarshalTypeError{
		Value:  "string",
		Type:   reflect.TypeOf(42),
		Offset: 26,
		Struct: "Customer",
		Field:  "age",
	}

	_, customer, k := initializeHandlersTest(t)
	inputReader := strings.NewReader(input)
	req := httptest.NewRequest("GET", "/dummy", inputReader)
	r := request.NewHTTPRequest(req)
	context2 := gofr.NewContext(nil, r, k)

	_, err := customer.Create(context2)
	if !reflect.DeepEqual(err, expErr) {
		t.Errorf("Testcase Failed\tExpected %v\tGot %v\n", expErr, err)
	}
}

func TestCustomer_Create(t *testing.T) {
	testCases := []struct {
		customer         string
		expectedResponse string
		err              error
	}{
		{`{"name":"Pirlo","age":42,"city":"Turin"}`, "New Customer Added!!", nil},
		{`{"name":"Pirlo","age":42,"city":"Turin"}`, "", errors2.New("test error")},
	}

	customerStore, customer, k := initializeHandlersTest(t)

	for i, tc := range testCases {
		input := strings.NewReader(tc.customer)
		req := httptest.NewRequest("GET", "/dummy", input)
		r := request.NewHTTPRequest(req)
		context2 := gofr.NewContext(nil, r, k)

		customerStore.EXPECT().Create(gomock.Any(), gomock.Any()).Return(tc.err)

		_, err := customer.Create(context2)
		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("Testcase[%v] Failed\tExpected %v\tGot %v\n", i, tc.err, err)
		}
	}
}

func TestCustomer_Delete(t *testing.T) {
	testCases := []struct {
		queryParams         string
		expectedDeleteCount int
		expectedResponse    string
		err                 error
	}{
		{"name=1", 0, "", errors.InvalidParam{Param: []string{"name"}}},
		{"name=Tim", 2, "2 Customers deleted!", nil},
		{"name=Thomas", 1, "1 Customers deleted!", nil},
	}

	customerStore, customer, k := initializeHandlersTest(t)

	for i, tc := range testCases {
		req := httptest.NewRequest("GET", "/customer?"+tc.queryParams, nil)
		r := request.NewHTTPRequest(req)
		context2 := gofr.NewContext(nil, r, k)

		customerStore.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(tc.expectedDeleteCount, tc.err).Times(1)

		_, err := customer.Delete(context2)

		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("Testcase[%v] failed\nExpected: %v \nGot: %v", i, tc.err, err)
		}
	}
}
