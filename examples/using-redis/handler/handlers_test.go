package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/metrics"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
)

type mockStore struct{}

const redisKey = "someKey"

func (m mockStore) Get(ctx *gofr.Context, key string) (string, error) {
	switch key {
	case redisKey:
		return "someValue", nil
	case "errorKey":
		return "", mockStore{}
	default:
		return "", mockStore{}
	}
}

func (m mockStore) Set(ctx *gofr.Context, key, value string, expiration time.Duration) error {
	if key == redisKey && value == "someValue" {
		return mockStore{}
	}

	return nil
}

func (m mockStore) Delete(ctx *gofr.Context, key string) error {
	switch key {
	case redisKey:
		return nil
	case "errorKey":
		return mockStore{}
	default:
		return nil
	}
}

func (m mockStore) Error() string {
	return "some mocked error"
}

func TestModel_GetKey(t *testing.T) {
	m := New(mockStore{})

	app := gofr.New()

	tests := []struct {
		desc string
		key  string
		resp interface{}
		err  error
	}{
		{"get with key", redisKey, "someValue", nil},
		{"get with empty key", "", nil, errors.MissingParam{Param: []string{"key"}}},
		{"get with error key", "errorKey", nil, mockStore{}},
	}

	for i, tc := range tests {
		r := httptest.NewRequest(http.MethodGet, "http://dummy", nil)

		req := request.NewHTTPRequest(r)
		ctx := gofr.NewContext(nil, req, app)

		if tc.key != "" {
			ctx.SetPathParams(map[string]string{
				"key": tc.key,
			})
		}

		_, err := m.GetKey(ctx)

		assert.Equal(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

func TestModel_DeleteKey(t *testing.T) {
	m := New(mockStore{})

	app := gofr.New()

	tests := []struct {
		desc string
		key  string
		err  error
	}{
		{"delete success", redisKey, nil},
		{"delete with empty key", "", errors.MissingParam{Param: []string{"key"}}},
		{"delete with error key", "errorKey", deleteErr{}},
	}

	for i, tc := range tests {
		r := httptest.NewRequest(http.MethodDelete, "http://dummy", nil)

		req := request.NewHTTPRequest(r)
		ctx := gofr.NewContext(nil, req, app)

		if tc.key != "" {
			ctx.SetPathParams(map[string]string{
				"key": tc.key,
			})
		}

		_, err := m.DeleteKey(ctx)

		assert.Equal(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

func TestModel_SetKey(t *testing.T) {
	m := New(mockStore{})

	app := gofr.New()
	mockMetric := metrics.NewMockMetric(gomock.NewController(t))
	app.Metric = mockMetric

	mockMetric.EXPECT().SetGauge(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockMetric.EXPECT().IncCounter(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	tests := []struct {
		desc string
		body []byte
		err  error
	}{
		{"set key with invalid body", []byte(`{`), invalidBodyErr{}},
		{"set key with invalid input", []byte(`{"someKey":"someValue"}`), invalidInputErr{}},
		{"set key success", []byte(`{"someKey123": "123"}`), nil},
	}

	for i, tc := range tests {
		r := httptest.NewRequest(http.MethodPost, "http://dummy", bytes.NewReader(tc.body))

		req := request.NewHTTPRequest(r)
		ctx := gofr.NewContext(nil, req, app)

		_, err := m.SetKey(ctx)

		assert.Equal(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

func TestSetKey_SetGaugeError(t *testing.T) {
	app := gofr.New()
	m := New(mockStore{})

	r := httptest.NewRequest(http.MethodPost, "http://dummy", nil)

	req := request.NewHTTPRequest(r)
	ctx := gofr.NewContext(nil, req, app)
	mockMetric := metrics.NewMockMetric(gomock.NewController(t))
	ctx.Metric = mockMetric

	expErr := errors.Error("error case")
	mockMetric.EXPECT().SetGauge(gomock.Any(), gomock.Any()).Return(expErr)

	_, err := m.SetKey(ctx)
	assert.Equal(t, expErr, err)
}

func TestSetKey_InvalidBodyCounterError(t *testing.T) {
	app := gofr.New()
	m := New(mockStore{})
	r := httptest.NewRequest(http.MethodPost, "http://dummy", bytes.NewReader([]byte(`{`)))
	req := request.NewHTTPRequest(r)
	ctx := gofr.NewContext(nil, req, app)
	mockMetric := metrics.NewMockMetric(gomock.NewController(t))
	ctx.Metric = mockMetric

	mockMetric.EXPECT().SetGauge(gomock.Any(), gomock.Any()).Return(nil)

	expErr := errors.Error("error case")
	mockMetric.EXPECT().IncCounter(gomock.Any()).Return(expErr)

	_, err := m.SetKey(ctx)
	assert.Equal(t, expErr, err)
}

func TestSetKey_IncCounterError(t *testing.T) {
	tcs := []struct {
		desc string
		body []byte
	}{
		{"invalid body", []byte(`{"`)},
		{"error key", []byte(`{"someKey":"someValue"}`)},
		{"valid key", []byte(`{"someKey1":"someValue1"}`)},
	}

	app := gofr.New()
	m := New(mockStore{})
	mockMetric := metrics.NewMockMetric(gomock.NewController(t))
	app.Metric = mockMetric
	expErr := errors.Error("error case")

	mockMetric.EXPECT().SetGauge(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockMetric.EXPECT().IncCounter(gomock.Any()).Return(nil)
	mockMetric.EXPECT().IncCounter(gomock.Any(), gomock.Any()).Return(expErr).AnyTimes()

	for i, tc := range tcs {
		r := httptest.NewRequest(http.MethodPost, "http://dummy", bytes.NewReader(tc.body))
		req := request.NewHTTPRequest(r)
		ctx := gofr.NewContext(nil, req, app)

		_, err := m.SetKey(ctx)
		assert.Equal(t, expErr, err, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

func TestDeleteErr_Error(t *testing.T) {
	var d deleteErr

	expected := "error: failed to delete"
	got := d.Error()

	if got != expected {
		t.Errorf("FAILED, expected: %v, got: %v", expected, got)
	}
}

func TestInvalidInputErr_Error(t *testing.T) {
	var i invalidInputErr

	expected := "error: invalid input"
	got := i.Error()

	if got != expected {
		t.Errorf("FAILED, expected: %v, got: %v", expected, got)
	}
}

func TestInvalidBodyErr_Error(t *testing.T) {
	var i invalidBodyErr

	expected := "error: invalid body"
	got := i.Error()

	if got != expected {
		t.Errorf("FAILED, expected: %v, got: %v", expected, got)
	}
}
