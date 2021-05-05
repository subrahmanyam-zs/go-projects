package main

import (
	"bytes"
	"net/http"
	"testing"
	"time"

	"github.com/zopsmart/gofr/pkg/gofr/request"
)

func TestServerRun(t *testing.T) {
	go main()
	time.Sleep(3 * time.Second)

	tcs := []struct {
		id                 int
		method             string
		endpoint           string
		expectedStatusCode int
		body               []byte
	}{
		{1, "GET", "http://localhost:9101", 101, nil},
		{2, "POST", "http://localhost:9101/ws", 405, nil},
		{3, "GET", "http://localhost:9101/ws", 101, nil},
	}

	for _, tc := range tcs {
		req, _ := request.NewMock(tc.method, tc.endpoint, bytes.NewBuffer(tc.body))
		req.Header.Set("Connection", "upgrade")
		req.Header.Set("Upgrade", "websocket")
		req.Header.Set("Sec-Websocket-Version", "13")
		req.Header.Set("Sec-WebSocket-Key", "wehkjeh21-sdjk210-wsknb")

		c := http.Client{}

		resp, _ := c.Do(req)
		if resp == nil {
			t.Errorf("Test %v: Failed \t got nil response", tc.id)
		}

		if resp != nil && resp.StatusCode != tc.expectedStatusCode {
			t.Errorf("Test %v: Failed.\tExpected %v\tGot %v\n", tc.id, tc.expectedStatusCode, resp.StatusCode)
		}
	}
}
