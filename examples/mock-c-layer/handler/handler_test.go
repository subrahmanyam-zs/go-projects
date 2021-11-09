package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"developer.zopsmart.com/go/gofr/examples/mock-c-layer/models"
	"developer.zopsmart.com/go/gofr/examples/mock-c-layer/store/brand"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"

	"github.com/stretchr/testify/assert"
)

func TestBrand_Get(t *testing.T) {
	tests := []struct {
		desc string
		id   string
		resp []models.Brand
		err  error
	}{
		{"get fail case", "4", nil, errors.Error("core error")},
		{"get single entity", "1", []models.Brand{{ID: 1, Name: "brand 1"}}, nil},
		{"get multiple entity", "2", []models.Brand{{ID: 1, Name: "brand 1"}, {ID: 2, Name: "brand 2"}}, nil},
	}

	store := brand.New()
	h := New(store)

	app := gofr.New()

	for i, tc := range tests {
		req := httptest.NewRequest(http.MethodGet, "/dummy?id="+tc.id, nil)
		r := request.NewHTTPRequest(req)
		ctx := gofr.NewContext(nil, r, app)
		data, err := h.Get(ctx)
		jsonData, _ := json.Marshal(data)

		assert.Equal(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)

		var b []models.Brand
		_ = json.Unmarshal(jsonData, &b)

		assert.Equal(t, tc.resp, b, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

func TestBrand_Create(t *testing.T) {
	tests := []struct {
		desc    string
		request []byte
		resp    interface{}
		err     error
	}{
		{"create empty brand", []byte(`{}`), models.Brand{}, nil},
		{"create success", []byte(`{"name":"Model 1"}`), models.Brand{Name: "Model 1"}, nil},
		{"create fail", []byte(`{"name":"brand 3"}`), nil, errors.Error("core error")},
	}

	store := brand.New()
	h := New(store)

	app := gofr.New()

	for i, tc := range tests {
		req := httptest.NewRequest(http.MethodGet, "/dummy", bytes.NewBuffer(tc.request))
		r := request.NewHTTPRequest(req)
		ctx := gofr.NewContext(nil, r, app)
		resp, err := h.Create(ctx)

		assert.Equal(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)

		var b models.Brand

		body, _ := json.Marshal(resp)

		_ = json.Unmarshal(body, &b)

		assert.Equal(t, tc.resp, resp, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

func TestBrand_CreateError(t *testing.T) {
	store := brand.New()
	h := New(store)

	app := gofr.New()
	expErr := errors.InvalidParam{Param: []string{"request body"}}
	body := []byte(`{"id":"1"}`)

	req := httptest.NewRequest(http.MethodGet, "/dummy", bytes.NewBuffer(body))
	r := request.NewHTTPRequest(req)

	ctx := gofr.NewContext(nil, r, app)
	_, err := h.Create(ctx)

	assert.Equal(t, expErr, err, "TEST, failed.\n")
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.Error("test error")
}

func TestBrand_CreateErrorBody(t *testing.T) {
	store := brand.New()
	h := New(store)

	app := gofr.New()
	expErr := errors.InvalidParam{Param: []string{"request body"}}
	req := httptest.NewRequest(http.MethodGet, "/dummy", errReader(0))
	r := request.NewHTTPRequest(req)
	ctx := gofr.NewContext(nil, r, app)
	_, err := h.Create(ctx)

	assert.Equal(t, expErr, err, "TEST, failed.\n")
}
