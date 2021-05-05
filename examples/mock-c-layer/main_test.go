package main

import (
	"bytes"
	"net/http"
	"testing"
	"time"

	"github.com/zopsmart/gofr/pkg/gofr/request"
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
		{"GET", "brand?id=1", 200, nil},
		{"POST", "brand", 201, []byte(`{"name":"brand 1"}`)},

		{"GET", "brand", 500, nil},
		{"POST", "brand", 500, []byte(`{"name":"brand 3"}`)},

		{"GET", "unknown", 404, nil},
		{"GET", "brand/id", 404, nil},

		{"PUT", "brand", 404, nil},
		{"DELETE", "brand", 404, nil},
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
