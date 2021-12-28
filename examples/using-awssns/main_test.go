//go:build !integration

package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
)

func TestIntegration(t *testing.T) {
	go main()
	time.Sleep(3 * time.Second)

	tests := []struct {
		method             string
		endpoint           string
		expectedStatusCode int
		body               []byte
	}{
		{http.MethodPost, "publish", http.StatusCreated, []byte(`{"name": "GOFR", "message":  "hi"}`)},
		{http.MethodGet, "subscribe", http.StatusOK, nil},
	}

	for i, tc := range tests {
		req, _ := request.NewMock(tc.method, "http://localhost:8080/"+tc.endpoint, http.NoBody)
		c := http.Client{}

		resp, err := c.Do(req)
		if err != nil {
			t.Errorf("TEST %v: error while making request err, %v", i+1, err)
			continue
		}

		assert.Equal(t, tc.expectedStatusCode, resp.StatusCode, "Test %v: Failed.\n", i+1)

		_ = resp.Body.Close()
	}
}
