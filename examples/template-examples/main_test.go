package main

import (
	"bytes"
	"net/http"
	"testing"
	"time"

	"github.com/zopsmart/gofr/pkg/gofr/request"
)

func TestServerIntegration(t *testing.T) {
	go main()
	time.Sleep(3 * time.Second)

	tcs := []struct {
		method              string
		endpoint            string
		expectedStatusCode  int
		body                []byte
		expectedContentType string
	}{
		{"GET", "test", 200, nil, "text/html"},
		{"GET", "test2", 404, nil, "application/json"},
		{"GET", "image", 200, nil, "image/png"},
	}

	for index, tc := range tcs {
		req, _ := request.NewMock(tc.method, "http://localhost:8000/"+tc.endpoint, bytes.NewBuffer(tc.body))
		c := http.Client{}

		resp, _ := c.Do(req)
		if resp != nil && resp.StatusCode != tc.expectedStatusCode {
			t.Errorf("Testcase[%v] Failed.\tExpected %v\tGot %v\n", index, tc.expectedStatusCode, resp.StatusCode)
		}

		if resp != nil && resp.Header.Get("Content-type") != tc.expectedContentType {
			t.Errorf("Testcase[%v] Failed.\tExpected %v\tGot %v\n", index, tc.expectedContentType, resp.Header.Get("Content-type"))
		}

		resp.Body.Close()
	}
}
