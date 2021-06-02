package main

import (
	"bytes"
	"net/http"
	"testing"
	"time"
)

func TestRoutes(t *testing.T) {
	go main()

	time.Sleep(time.Second * 15)

	testcases := []struct {
		method             string
		endpoint           string
		expectedStatusCode int
		body               []byte
	}{
		{"GET", "unknown", 404, nil},
		{"POST", "unknown", 404, nil},
		{"OPOOSOS", "unknown", 404, nil},
	}

	for _, tc := range testcases {
		req, _ := http.NewRequest(tc.method, "http://localhost:8001/"+tc.endpoint, bytes.NewBuffer(tc.body))
		c := http.Client{}

		resp, err := c.Do(req)
		if err != nil {
			t.Errorf("got error: %v", err)
		}

		if resp != nil && resp.StatusCode != tc.expectedStatusCode {
			t.Errorf("Failed.\tExpected %v\tGot %v\n", tc.expectedStatusCode, resp.StatusCode)
		}

		_ = resp.Body.Close()
	}
}
