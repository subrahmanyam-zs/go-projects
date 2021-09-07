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
		{method: http.MethodPost, endpoint: "publish", expectedStatusCode: 201, body: []byte(`{"name": "GOFR", "message":  "hi"}`)},
		{method: http.MethodGet, endpoint: "subscribe", expectedStatusCode: 200, body: nil},
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
				assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)
				resp.Body.Close()
			}
		})
	}
}
