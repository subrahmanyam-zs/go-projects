package main

import (
	"bytes"
	"net/http"
	"testing"
	"time"

	"github.com/zopsmart/gofr/pkg/gofr/request"
)

func TestServerValidation(t *testing.T) {
	go main()
	time.Sleep(3 * time.Second)

	testcases := []struct {
		method             string
		endpoint           string
		expectedStatusCode int
		body               []byte
	}{
		{"POST", "phone", 201, []byte(`{"phone":"+912123456789098", "email": "c.r@yahoo.com"}`)},
		{"POST", "phone", 500, nil},
		{"POST", "phone2", 404, nil},
		{"GET", "phone", 404, nil},
	}

	for index, tc := range testcases {
		req, _ := request.NewMock(tc.method, "http://localhost:9010/"+tc.endpoint, bytes.NewBuffer(tc.body))
		c := http.Client{}

		resp, _ := c.Do(req)
		if resp != nil && resp.StatusCode != tc.expectedStatusCode {
			t.Errorf("Test Case: %v \tFailed.\tExpected %v\tGot %v\n", index+1, tc.expectedStatusCode, resp.StatusCode)
		}
	}
}
