package main

import (
	"bytes"
	"net/http"
	"testing"
	"time"

	"github.int.mcafee.com/mcafee/cnsr-gofr-csp-auth/generator"
)

func TestIntegration(t *testing.T) {
	go main()
	time.Sleep(time.Second * 5)

	csp, err := generator.New(&generator.Option{
		AppKey:     "mock-app-key",
		SharedKey:  "mock-shared-key",
		TimeFormat: time.RFC3339,
	})
	if err != nil {
		t.Errorf("error while creating instance of generator for csp security, %v", err)
	}

	tests := []struct {
		desc          string
		method        string
		endPoint      string
		body          []byte
		expStatusCode int
	}{
		{"valid request", http.MethodGet, "hello", nil, http.StatusOK},
		{"valid request", http.MethodPost, "greet", []byte(`{"name":"test"}`), http.StatusCreated},
		{"unknown route", http.MethodDelete, "invalidRoute", nil, http.StatusNotFound},
	}

	for i, test := range tests {
		req, _ := http.NewRequest(test.method, "http://localhost:4000/"+test.endPoint, bytes.NewBuffer(test.body))
		headers := csp.GetCSPHeaders(req)

		// set csp security headers
		for key, val := range headers {
			req.Header.Add(key, val)
		}

		client := http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("error while making request, %v", err)
			continue
		}

		if resp.StatusCode != test.expStatusCode {
			t.Errorf("TEST[%v] %v\nExpected status code %v, got %v", i, test.desc, test.expStatusCode, resp.StatusCode)
		}

		_ = resp.Body.Close()
	}
}
