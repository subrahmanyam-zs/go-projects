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

	tcs := []struct {
		method             string
		endpoint           string
		expectedStatusCode int
		body               []byte
	}{
		{"GET", "config/key123", 500, nil},
		{"POST", "config", 201, []byte(`{}`)},
		{"DELETE", "config/key123", 204, nil},
	}

	for _, tc := range tcs {
		req, _ := request.NewMock(tc.method, "http://localhost:9091/"+tc.endpoint, bytes.NewBuffer(tc.body))
		c := http.Client{}

		resp, _ := c.Do(req)

		if resp != nil && resp.StatusCode != tc.expectedStatusCode {
			t.Errorf("Failed.\tExpected %v\tGot %v\n", tc.expectedStatusCode, resp.StatusCode)
		}

		resp.Body.Close()
	}
}
