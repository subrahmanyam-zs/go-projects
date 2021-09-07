package handler

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
)

func TestHelloWorldHandler(t *testing.T) {
	c := gofr.NewContext(nil, nil, nil)

	resp, err := HelloWorld(c)
	if err != nil {
		t.Errorf("FAILED, got error: %v", err)
	}

	expected := "Hello World!"
	got := fmt.Sprintf("%v", resp)

	if got != expected {
		t.Errorf("FAILED, Expected: %v, Got: %v", expected, got)
	}
}

func TestHelloNameHandler(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{
			name:     "SomeName",
			expected: "Hello SomeName",
		},
		{
			name:     "Firstname Lastname",
			expected: "Hello Firstname Lastname",
		},
	}

	for _, test := range tests {
		r := httptest.NewRequest("GET", "http://dummy/hello?name="+url.QueryEscape(test.name), nil)
		req := request.NewHTTPRequest(r)
		c := gofr.NewContext(nil, req, nil)

		resp, err := HelloName(c)
		if err != nil {
			t.Errorf("FAILED, got error: %v", err)
		}

		got := fmt.Sprintf("%v", resp)
		if got != test.expected {
			t.Errorf("FAILED, Expected: %v, Got: %v", test.expected, got)
		}
	}
}

func TestJSONHandler(t *testing.T) {
	c := gofr.NewContext(nil, nil, nil)

	res, err := JSONHandler(c)
	if err != nil {
		t.Errorf("FAILED, got error: %v", err)
	}

	expected := resp{
		Name:    "Vikash",
		Company: "ZopSmart",
	}

	var got resp

	resBytes, _ := json.Marshal(res)

	if err := json.Unmarshal(resBytes, &got); err != nil {
		t.Errorf("FAILED, got error: %v", err)
	}

	if !reflect.DeepEqual(expected, got) {
		t.Errorf("FAILED, Expected: %v, Got: %v", expected, got)
	}
}

func TestUserHandler(t *testing.T) {
	tests := []struct {
		name string
		resp interface{}
		err  error
	}{
		{"Vikash", resp{Name: "Vikash", Company: "ZopSmart"}, nil},
		{"ABC", nil, errors.EntityNotFound{Entity: "user", ID: "ABC"}},
	}

	for _, tc := range tests {
		r := httptest.NewRequest("GET", "http://dummy", nil)
		req := request.NewHTTPRequest(r)

		c := gofr.NewContext(nil, req, nil)
		c.SetPathParams(map[string]string{"name": tc.name})

		resp, err := UserHandler(c)

		assert.Equal(t, tc.err, err)

		assert.Equal(t, tc.resp, resp)
	}
}

func TestErrorHandler(t *testing.T) {
	c := gofr.NewContext(nil, nil, nil)

	res, err := ErrorHandler(c)
	if res != nil {
		t.Errorf("FAILED, expected nil, got: %v", res)
	}

	exp := &errors.Response{
		StatusCode: 500,
		Code:       "UNKNOWN_ERROR",
		Reason:     "unknown error occurred",
	}

	if !assert.Equal(t, exp, err) {
		t.Errorf("FAILED, exp: %v, got: %v", exp, err)
	}
}

func TestHelloLogHandler(t *testing.T) {
	r := httptest.NewRequest("GET", "http://dummy/log", nil)
	req := request.NewHTTPRequest(r)
	c := gofr.NewContext(nil, req, nil)

	res, err := HelloLogHandler(c)
	if res != "Logging OK" {
		t.Errorf("Logging Failed due to : %v", err)
	}
}
