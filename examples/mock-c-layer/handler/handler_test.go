package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"developer.zopsmart.com/go/gofr/examples/mock-c-layer/store"
	"developer.zopsmart.com/go/gofr/examples/mock-c-layer/store/brand"
	errors2 "developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
)

func TestBrand_Get(t *testing.T) {
	//nolint:govet  // table tests
	testCases := []struct {
		reqID        string
		expectedResp []store.Model
		err          error
	}{
		{"4", nil, errors.New("core error")},

		{"1", []store.Model{{1, "brand 1"}}, nil},
		{"2", []store.Model{{1, "brand 1"}, {2, "brand 2"}}, nil},
	}

	brandCore := brand.New()
	consumer := New(brandCore)

	k := gofr.New()

	for _, tc := range testCases {
		req := httptest.NewRequest(http.MethodGet, "/dummy?id="+tc.reqID, nil)
		r := request.NewHTTPRequest(req)
		context := gofr.NewContext(nil, r, k)
		data, err := consumer.Get(context)
		jsonData, _ := json.Marshal(data)

		if !reflect.DeepEqual(tc.err, err) {
			t.Errorf("Failed\tExpected %v\tGot %v\n", tc.err, err)
		}

		var b []store.Model
		_ = json.Unmarshal(jsonData, &b)

		if !reflect.DeepEqual(b, tc.expectedResp) {
			t.Errorf("Retrieval from core layer failed")
		}
	}
}
func TestBrand_Create(t *testing.T) {
	testCases := []struct {
		request          []byte
		expectedResponse store.Model
		err              error
	}{
		{[]byte(`{}`), store.Model{}, nil},
		{[]byte(`{"name":"Model 1"}`), store.Model{Name: "Model 1"}, nil},
		{[]byte(`{"name":"brand 3"}`), store.Model{}, errors.New("core error")},
	}

	brandCore := brand.New()
	consumer := New(brandCore)

	k := gofr.New()

	for _, tc := range testCases {
		req := httptest.NewRequest(http.MethodGet, "/dummy", bytes.NewBuffer(tc.request))
		r := request.NewHTTPRequest(req)
		context := gofr.NewContext(nil, r, k)
		resp, err := consumer.Create(context)

		if !reflect.DeepEqual(tc.err, err) {
			t.Errorf("Failed\tExpected %v\tGot %v\n", tc.err, err)
		}

		var b store.Model

		body, _ := json.Marshal(resp)

		_ = json.Unmarshal(body, &b)

		if !reflect.DeepEqual(b, tc.expectedResponse) {
			t.Errorf("Retrieval from core layer failed")
		}
	}
}

func TestBrand_CreateError(t *testing.T) {
	brandCore := brand.New()
	consumer := New(brandCore)

	k := gofr.New()
	expectedError := errors2.InvalidParam{Param: []string{"request body"}}
	body := []byte(`{"id":"1"}`)

	req := httptest.NewRequest(http.MethodGet, "/dummy", bytes.NewBuffer(body))
	r := request.NewHTTPRequest(req)

	context := gofr.NewContext(nil, r, k)
	_, err := consumer.Create(context)

	if err.Error() != expectedError.Error() {
		t.Errorf("Failed\tExpected %v\tGot %v\n", expectedError, err)
	}
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func TestBrand_CreateErrorBody(t *testing.T) {
	brandCore := brand.New()
	consumer := New(brandCore)

	k := gofr.New()
	expectedError := errors2.InvalidParam{Param: []string{"request body"}}
	req := httptest.NewRequest(http.MethodGet, "/dummy", errReader(0))
	r := request.NewHTTPRequest(req)
	context := gofr.NewContext(nil, r, k)
	_, err := consumer.Create(context)

	if err.Error() != expectedError.Error() {
		t.Errorf("Failed\tExpected %v\tGot %v\n", expectedError, err)
	}
}
