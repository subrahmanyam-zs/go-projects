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
	time.Sleep(3 * time.Second)

	tcs := []struct {
		method             string
		endpoint           string
		expectedStatusCode int
		body               []byte
	}{
		{"GET", "hello-world", 200, nil},
		{"GET", "hello?name=random", 200, nil},
		{"GET", "json", 200, nil},
		{"GET", "error", 500, nil},
		{"GET", "swagger", 200, nil},
	}

	for _, tc := range tcs {
		req, _ := request.NewMock(tc.method, "http://localhost:9000/"+tc.endpoint, bytes.NewBuffer(tc.body))

		c := http.Client{}
		resp, _ := c.Do(req)

		if resp != nil && resp.StatusCode != tc.expectedStatusCode {
			t.Errorf("Failed.\tExpected %v\tGot %v\n", tc.expectedStatusCode, resp.StatusCode)
		}

		resp.Body.Close()
	}
}
