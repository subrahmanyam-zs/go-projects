package handler

import (
	"bytes"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zopsmart/gofr/pkg/errors"
	"github.com/zopsmart/gofr/pkg/gofr"
	"github.com/zopsmart/gofr/pkg/gofr/request"
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

func TestRedisModel_GetKey(t *testing.T) {
	m := New(mockStore{})
	k := gofr.New()

	tests := []struct {
		key         string
		value       string
		expectedErr error
	}{
		{redisKey, "someValue", nil},
		{"", "emptyKeyValue", errors.MissingParam{Param: []string{"key"}}},
		{"errorKey", "errorKeyValue", mockStore{}},
	}

	for _, tc := range tests {
		r := httptest.NewRequest("GET", "http://dummy", nil)
		req := request.NewHTTPRequest(r)
		c := gofr.NewContext(nil, req, k)

		if tc.key != "" {
			c.SetPathParams(map[string]string{
				"key": tc.key,
			})
		}

		_, gotErr := m.GetKey(c)
		assert.Equal(t, tc.expectedErr, gotErr)
	}
}

func TestRedisModel_SetKey(t *testing.T) {
	m := New(mockStore{})
	k := gofr.New()

	tests := []struct {
		body        []byte
		expectedErr error
	}{
		{[]byte(`{`), invalidBodyErr},
		{[]byte(`{"someKey":"someValue"}`), invalidInputErr},
		{[]byte(`{"someKey123": "123"}`), nil},
	}

	for _, tc := range tests {
		r := httptest.NewRequest("POST", "http://dummy", bytes.NewReader(tc.body))
		req := request.NewHTTPRequest(r)
		c := gofr.NewContext(nil, req, k)

		_, gotErr := m.SetKey(c)
		assert.Equal(t, tc.expectedErr, gotErr)
	}
}

func (m mockStore) Error() string {
	return "some mocked error"
}

func TestRedisInvalidInputErr_Error(t *testing.T) {
	err := constError("error: invalid input")
	expected := "error: invalid input"
	got := err.Error()

	assert.Equal(t, expected, got)
}

func TestRedisInvalidBodyErr_Error(t *testing.T) {
	err := constError("error: invalid body")
	expected := "error: invalid body"
	got := err.Error()

	assert.Equal(t, expected, got)
}
