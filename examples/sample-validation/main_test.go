package main

import (
	"bytes"
	"net/http"
	"testing"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
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
		{http.MethodPost, "phone", 201, []byte(`{"phone":"+912123456789098", "email": "c.r@yahoo.com"}`)},
		{http.MethodPost, "phone", 500, nil},
		{http.MethodPost, "phone2", 404, nil},
		{http.MethodGet, "phone", 404, nil},
	}

	for index, tc := range testcases {
		req, _ := request.NewMock(tc.method, "http://localhost:9010/"+tc.endpoint, bytes.NewBuffer(tc.body))
		c := http.Client{}

		resp, err := c.Do(req)
		if err != nil {
			t.Errorf("error on making request , %v", err)
		}

		if resp != nil && resp.StatusCode != tc.expectedStatusCode {
			t.Errorf("Test Case: %v \tFailed.\tExpected %v\tGot %v\n", index+1, tc.expectedStatusCode, resp.StatusCode)
		}

		_ = resp.Body.Close()
	}
}
