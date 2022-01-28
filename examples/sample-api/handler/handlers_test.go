package handler

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
)

func TestHelloWorldHandler(t *testing.T) {
	ctx := gofr.NewContext(nil, nil, nil)

	resp, err := HelloWorld(ctx)
	if err != nil {
		t.Errorf("FAILED, Expected: %v, Got: %v", nil, err)
	}

	expected := "Hello World!"
	got := fmt.Sprintf("%v", resp)

	if got != expected {
		t.Errorf("FAILED, Expected: %v, Got: %v", expected, got)
	}
}

func TestHelloNameHandler(t *testing.T) {
	tests := []struct {
		desc string
		name string
		resp string
	}{
		{"hello without lastname", "SomeName", "Hello SomeName"},
		{"hello with lastname", "Firstname Lastname", "Hello Firstname Lastname"},
	}

	for i, tc := range tests {
		r := httptest.NewRequest("GET", "http://dummy/hello?name="+url.QueryEscape(tc.name), nil)
		req := request.NewHTTPRequest(r)
		ctx := gofr.NewContext(nil, req, nil)

		resp, err := HelloName(ctx)

		assert.Equal(t, nil, err, "TEST[%d], failed.\n%s", i, tc.desc)

		assert.Equal(t, tc.resp, resp, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

func TestJSONHandler(t *testing.T) {
	ctx := gofr.NewContext(nil, nil, nil)

	res, err := JSONHandler(ctx)
	if err != nil {
		t.Errorf("FAILED, got error: %v", err)
	}

	expected := resp{Name: "Vikash", Company: "ZopSmart"}

	var got resp

	resBytes, _ := json.Marshal(res)

	if err := json.Unmarshal(resBytes, &got); err != nil {
		t.Errorf("FAILED, got error: %v", err)
	}

	assert.Equal(t, expected, got)
}

func TestUserHandler(t *testing.T) {
	tests := []struct {
		desc string
		name string
		resp interface{}
		err  error
	}{
		{"UserHandler success", "Vikash", resp{Name: "Vikash", Company: "ZopSmart"}, nil},
		{"UserHandler fail", "ABC", nil, errors.EntityNotFound{Entity: "user", ID: "ABC"}},
	}

	for i, tc := range tests {
		r := httptest.NewRequest("GET", "http://dummy", nil)
		req := request.NewHTTPRequest(r)

		ctx := gofr.NewContext(nil, req, nil)
		ctx.SetPathParams(map[string]string{"name": tc.name})

		resp, err := UserHandler(ctx)

		assert.Equal(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)

		assert.Equal(t, tc.resp, resp, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

func TestErrorHandler(t *testing.T) {
	ctx := gofr.NewContext(nil, nil, nil)

	res, err := ErrorHandler(ctx)
	if res != nil {
		t.Errorf("FAILED, expected nil, got: %v", res)
	}

	exp := &errors.Response{
		StatusCode: 500,
		Code:       "UNKNOWN_ERROR",
		Reason:     "unknown error occurred",
	}

	assert.Equal(t, exp, err)
}

func TestHelloLogHandler(t *testing.T) {
	r := httptest.NewRequest("GET", "http://dummy/log", nil)
	req := request.NewHTTPRequest(r)
	ctx := gofr.NewContext(nil, req, nil)

	res, err := HelloLogHandler(ctx)
	if res != "Logging OK" {
		t.Errorf("Logging Failed due to : %v", err)
	}
}
