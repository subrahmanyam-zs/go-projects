package main

import (
	"bytes"
	"crypto/tls"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/zopsmart/gofr/pkg/gofr/request"
)

func TestServerRun(t *testing.T) {
	os.Setenv("VALIDATE_HEADERS", "Custom-Header")
	defer os.Unsetenv("VALIDATE_HEADERS")

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
		headers            map[string]string
		body               []byte
	}{
		{1, "GET", "https://localhost:1443/hello-world", 200, nil, nil},
		{2, "GET", "https://localhost:1443/hello-world/", 200, nil, nil},
		{3, "POST", "https://localhost:1443/post", 201, nil, []byte(`{"Username":"username"}`)},
		{4, "POST", "https://localhost:1443/post/", 200, nil, []byte(`{"Username":"alreadyExist"}`)},
		// http will be redirected to https as redirect is set to true in https configuration
		{5, "GET", "http://localhost:9007/hello?name=random", 200, nil, nil},
		{6, "GET", "http://localhost:9007/multiple-errors", 500, nil, nil},
		{6, "GET", "http://localhost:9007/multiple-errors?id=1", 400, nil, nil},
		{7, "GET", "http://localhost:9007/.well-known/heartbeat", 200,
			map[string]string{"Content-Type": "application/json"}, nil},
		{7, "GET", "http://localhost:9007/hello", 400,
			map[string]string{"X-Zopsmart-Tenant": "food"}, nil},
		{8, "GET", "http://localhost:9007/error", 404, nil, nil},
	}

	for _, tc := range tcs {
		req, _ := request.NewMock(tc.method, tc.endpoint, bytes.NewBuffer(tc.body))
		c := http.Client{Transport: tr}

		if tc.headers == nil {
			req.Header.Set("Custom-Header", "test")
		}

		resp, _ := c.Do(req)

		if resp == nil {
			t.Errorf("Test %v: Failed \t got nil response", tc.id)
			return
		}

		if resp.StatusCode != tc.expectedStatusCode {
			t.Errorf("Test %v: Failed.\tExpected %v\tGot %v\n", tc.id, tc.expectedStatusCode, resp.StatusCode)
		}

		resp.Body.Close()
	}
}
