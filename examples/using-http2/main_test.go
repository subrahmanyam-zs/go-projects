package main

import (
	"bytes"
	"crypto/tls"
	"net/http"
	"testing"
	"time"

	"github.com/zopsmart/gofr/pkg/gofr/request"
)

func TestServerRun(t *testing.T) {
	go main()
	time.Sleep(3 * time.Second)

	//nolint: gosec, TLS InsecureSkipVerify set true.
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	tcs := []struct {
		id                 int
		method             string
		endpoint           string
		expectedStatusCode int
		body               []byte
	}{
		{1, "GET", "https://localhost:1449/home", 200, nil},
	}
	for _, tc := range tcs {
		req, _ := request.NewMock(tc.method, tc.endpoint, bytes.NewBuffer(tc.body))
		c := http.Client{Transport: tr}

		resp, _ := c.Do(req)

		if resp == nil {
			t.Errorf("Test %v: Failed \t got nil response", tc.id)
		}

		if resp != nil && resp.StatusCode != tc.expectedStatusCode {
			t.Errorf("Test %v: Failed.\tExpected %v\tGot %v\n", tc.id, tc.expectedStatusCode, resp.StatusCode)
		}

		resp.Body.Close()
	}
}
