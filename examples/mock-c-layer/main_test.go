package main

import (
	"bytes"
	"net/http"
	"testing"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
)

func TestIntegration(t *testing.T) {
	go main()
	time.Sleep(time.Second * 5)

	tests := []struct {
		desc       string
		method     string
		endpoint   string
		statusCode int
		body       []byte
	}{
		{"get success", http.MethodGet, "brand?id=1", http.StatusOK, nil},
		{"create success", http.MethodPost, "brand", http.StatusCreated, []byte(`{"name":"brand 1"}`)},

		{"get fail", http.MethodGet, "brand", http.StatusInternalServerError, nil},
		{"create fail", http.MethodPost, "brand", http.StatusInternalServerError, []byte(`{"name":"brand 3"}`)},

		{"get invalid route", http.MethodGet, "unknown", http.StatusNotFound, nil},
		{"get invalid endpoint", http.MethodGet, "brand/id", http.StatusNotFound, nil},

		{"unregistered update route", http.MethodPut, "brand", http.StatusMethodNotAllowed, nil},
		{"unregistered delete route", http.MethodDelete, "brand", http.StatusMethodNotAllowed, nil},
	}

	for i, tc := range tests {
		req, _ := request.NewMock(tc.method, "http://localhost:9090/"+tc.endpoint, bytes.NewBuffer(tc.body))
		c := http.Client{}

		resp, err := c.Do(req)
		if err != nil {
			t.Errorf("TEST[%v] Failed.\tHTTP request encountered Err: %v\n%s", i, err, tc.desc)
			continue
		}

		if resp.StatusCode != tc.statusCode {
			t.Errorf("TEST[%v] Failed.\tExpected %v\tGot %v\n%s", i, tc.statusCode, resp.StatusCode, tc.desc)
		}

		_ = resp.Body.Close()
	}
}
