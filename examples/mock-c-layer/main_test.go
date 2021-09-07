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

	testcases := []struct {
		method             string
		endpoint           string
		expectedStatusCode int
		body               []byte
	}{
		{http.MethodGet, "brand?id=1", 200, nil},
		{http.MethodPost, "brand", 201, []byte(`{"name":"brand 1"}`)},

		{http.MethodGet, "brand", 500, nil},
		{http.MethodPost, "brand", 500, []byte(`{"name":"brand 3"}`)},

		{http.MethodGet, "unknown", 404, nil},
		{http.MethodGet, "brand/id", 404, nil},

		{http.MethodPut, "brand", 404, nil},
		{http.MethodDelete, "brand", 404, nil},
	}

	for _, tc := range testcases {
		req, _ := request.NewMock(tc.method, "http://localhost:9090/"+tc.endpoint, bytes.NewBuffer(tc.body))
		c := http.Client{}

		resp, _ := c.Do(req)

		if resp != nil && resp.StatusCode != tc.expectedStatusCode {
			t.Errorf("Failed.\tExpected %v\tGot %v\n", tc.expectedStatusCode, resp.StatusCode)
		}

		resp.Body.Close()
	}
}
