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
		{"GET", "customer?name=Name", 200, nil},
		{"POST", "customer", 201, []byte(`{"name":"Robert"}`)},

		{"GET", "unknown", 404, nil},
		{"GET", "customer/id", 404, nil},

		{"PUT", "customer", 404, nil},
		{"DELETE", "customer?name=Robert", 204, nil},
	}

	for index, tc := range testcases {
		req, _ := request.NewMock(tc.method, "http://localhost:9097/"+tc.endpoint, bytes.NewBuffer(tc.body))
		c := http.Client{}

		resp, _ := c.Do(req)

		if resp != nil && resp.StatusCode != tc.expectedStatusCode {
			t.Errorf("Testcase[%v] Failed.\tExpected %v\tGot %v\n", index, tc.expectedStatusCode, resp.StatusCode)
		}

		resp.Body.Close()
	}
}
