//go:build !integration

package main

import (
	"fmt"
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
		{http.MethodPost, "http://localhost:8080/publish", http.StatusCreated, []byte(`{"name": "GOFR", "message":  "hi"}`)},
		{http.MethodGet, "http://localhost:8080/subscribe", http.StatusOK, nil},
	}

	for i, tc := range tests {
		tc := tc
		i := i
		t.Run(fmt.Sprintf("Test %v", i+1), func(t *testing.T) {
			req, _ := request.NewMock(tc.method, tc.endpoint, nil)
			c := http.Client{}

			resp, err := c.Do(req)
			if resp == nil || err != nil {
				t.Error(err)
			}

			if resp != nil {
				assert.Equal(t, tc.expectedStatusCode, resp.StatusCode, "Test %v: Failed.\n", i+1)

				err = resp.Body.Close()
				if err != nil {
					t.Error(err)
				}
			}
		})
	}
}
