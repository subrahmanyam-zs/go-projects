package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

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

	k := gofr.New()

	tests := []struct {
		key         string
		expectedErr error
		value       string
	}{
		{
			key:   redisKey,
			value: "someValue",
		},
		{
			key:         "",
			expectedErr: errors.MissingParam{Param: []string{"key"}},
		},
		{
			key:         "errorKey",
			expectedErr: mockStore{},
		},
	}

	for _, test := range tests {
		r := httptest.NewRequest("GET", "http://dummy", nil)

		req := request.NewHTTPRequest(r)
		c := gofr.NewContext(nil, req, k)

		if test.key != "" {
			c.SetPathParams(map[string]string{
				"key": test.key,
			})
		}

		_, gotErr := m.GetKey(c)

		if !reflect.DeepEqual(gotErr, test.expectedErr) {
			t.Errorf("FAILED, Expected: %v, Got: %v", test.expectedErr, gotErr)
		}
	}
}

func TestModel_DeleteKey(t *testing.T) {
	m := New(mockStore{})

	k := gofr.New()

	tests := []struct {
		key         string
		expectedErr error
	}{
		{
			key: redisKey,
		},
		{
			key:         "",
			expectedErr: errors.MissingParam{Param: []string{"key"}},
		},
		{
			key:         "errorKey",
			expectedErr: deleteErr{},
		},
	}

	for _, test := range tests {
		r := httptest.NewRequest("DELETE", "http://dummy", nil)

		req := request.NewHTTPRequest(r)
		c := gofr.NewContext(nil, req, k)

		if test.key != "" {
			c.SetPathParams(map[string]string{
				"key": test.key,
			})
		}

		_, gotErr := m.DeleteKey(c)

		if !reflect.DeepEqual(gotErr, test.expectedErr) {
			t.Errorf("FAILED, Expected: %v, Got: %v", test.expectedErr, gotErr)
		}
	}
}

func TestModel_SetKey(t *testing.T) {
	m := New(mockStore{})

	k := gofr.New()
	mockMetric := metrics.NewMockMetric(gomock.NewController(t))
	k.Metric = mockMetric

	mockMetric.EXPECT().SetGauge(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockMetric.EXPECT().IncCounter(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	tests := []struct {
		body        []byte
		expectedErr error
	}{
		{
			body:        []byte(`{`),
			expectedErr: invalidBodyErr{},
		},
		{
			body:        []byte(`{"someKey":"someValue"}`),
			expectedErr: invalidInputErr{},
		},
		{
			body:        []byte(`{"someKey123": "123"}`),
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		r := httptest.NewRequest("POST", "http://dummy", bytes.NewReader(test.body))

		req := request.NewHTTPRequest(r)
		c := gofr.NewContext(nil, req, k)

		_, gotErr := m.SetKey(c)

		if !reflect.DeepEqual(gotErr, test.expectedErr) {
			t.Errorf("FAILED, Expected: %v, Got: %v", test.expectedErr, gotErr)
		}
	}
}

func TestSetKey_SetGaugeError(t *testing.T) {
	k := gofr.New()
	m := New(mockStore{})

	r := httptest.NewRequest("POST", "http://dummy", nil)

	req := request.NewHTTPRequest(r)
	c := gofr.NewContext(nil, req, k)
	mockMetric := metrics.NewMockMetric(gomock.NewController(t))
	c.Metric = mockMetric

	expErr := errors.Error("error case")
	mockMetric.EXPECT().SetGauge(gomock.Any(), gomock.Any()).Return(expErr)

	_, err := m.SetKey(c)
	if !reflect.DeepEqual(err, expErr) {
		t.Errorf("expected error %v, got %v", expErr, err)
	}
}

func TestSetKey_InvalidBodyCounterError(t *testing.T) {
	k := gofr.New()
	m := New(mockStore{})
	r := httptest.NewRequest(http.MethodPost, "http://dummy", bytes.NewReader([]byte(`{`)))
	req := request.NewHTTPRequest(r)
	c := gofr.NewContext(nil, req, k)
	mockMetric := metrics.NewMockMetric(gomock.NewController(t))
	c.Metric = mockMetric

	mockMetric.EXPECT().SetGauge(gomock.Any(), gomock.Any()).Return(nil)

	expErr := errors.Error("error case")
	mockMetric.EXPECT().IncCounter(gomock.Any()).Return(expErr)

	_, err := m.SetKey(c)
	if !reflect.DeepEqual(err, expErr) {
		t.Errorf("expected error %v, got %v", expErr, err)
	}
}

func TestSetKey_IncCounterError(t *testing.T) {
	tcs := []struct {
		body []byte
	}{
		{[]byte(`{"`)},
		{[]byte(`{"someKey":"someValue"}`)},
		{[]byte(`{"someKey1":"someValue1"}`)},
	}

	k := gofr.New()
	m := New(mockStore{})
	mockMetric := metrics.NewMockMetric(gomock.NewController(t))
	k.Metric = mockMetric
	expErr := errors.Error("error case")
	mockMetric.EXPECT().SetGauge(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockMetric.EXPECT().IncCounter(gomock.Any()).Return(nil)
	mockMetric.EXPECT().IncCounter(gomock.Any(), gomock.Any()).Return(expErr).AnyTimes()

	for i, tc := range tcs {
		r := httptest.NewRequest("POST", "http://dummy", bytes.NewReader(tc.body))
		req := request.NewHTTPRequest(r)
		c := gofr.NewContext(nil, req, k)

		_, err := m.SetKey(c)
		if !reflect.DeepEqual(err, expErr) {
			t.Errorf("TESTCASE[%v] expected error %v, got %v", i, expErr, err)
		}
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
