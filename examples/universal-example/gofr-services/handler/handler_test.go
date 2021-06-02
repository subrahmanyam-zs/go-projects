package handler

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
	"developer.zopsmart.com/go/gofr/pkg/gofr/responder"
)

type mockService struct{}

func TestHandler_GetLog(t *testing.T) {
	tests := []struct {
		filter   string
		response interface{}
		err      error
	}{
		{"service=gofr-hello-api", "warn", nil},
		{"", nil, errors.MissingParam{Param: []string{"service"}}},
	}

	for i, tc := range tests {
		req := httptest.NewRequest("GET", "http://dummy?"+tc.filter, nil)
		c := gofr.NewContext(responder.NewContextualResponder(httptest.NewRecorder(), req), request.NewHTTPRequest(req), nil)

		h := New(mockService{})

		resp, err := h.Log(c)
		assert.Equal(t, tc.err, err, i)
		assert.Equal(t, tc.response, resp, i)
	}
}

func TestHandler_GetHello(t *testing.T) {
	tests := []struct {
		filter   string
		response interface{}
		err      error
	}{
		{"name=ZopSmart", "Hello ZopSmart", nil},
		{"", "Hello", nil},
	}

	for i, tc := range tests {
		req := httptest.NewRequest("GET", "http://dummy?"+tc.filter, nil)
		c := gofr.NewContext(responder.NewContextualResponder(httptest.NewRecorder(), req), request.NewHTTPRequest(req), nil)

		h := New(mockService{})

		resp, err := h.Hello(c)
		assert.Equal(t, tc.err, err, i)
		assert.Equal(t, tc.response, resp, i)
	}
}

func (m mockService) Log(ctx *gofr.Context, serviceName string) (string, error) {
	if serviceName != "" {
		return "warn", nil
	}

	return "", errors.MissingParam{Param: []string{"service"}}
}

func (m mockService) Hello(ctx *gofr.Context, name string) (string, error) {
	if name != "" {
		return "Hello " + name, nil
	}

	return "Hello", nil
}
