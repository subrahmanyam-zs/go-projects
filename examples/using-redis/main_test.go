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

	tcs := []struct {
		method             string
		endpoint           string
		expectedStatusCode int
		body               []byte
	}{
		{http.MethodGet, "config/key123", 500, nil},
		{http.MethodPost, "config", 201, []byte(`{}`)},
		{http.MethodDelete, "config/key123", 204, nil},
	}

	for _, tc := range tcs {
		req, _ := request.NewMock(tc.method, "http://localhost:9091/"+tc.endpoint, bytes.NewBuffer(tc.body))
		c := http.Client{}

		resp, err := c.Do(req)
		if resp == nil || err != nil {
			t.Error(err)
			continue
		}

		if resp.StatusCode != tc.expectedStatusCode {
			t.Errorf("Failed.\tExpected %v\tGot %v\n", tc.expectedStatusCode, resp.StatusCode)
		}

		_ = resp.Body.Close()
	}
}
