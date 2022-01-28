//go:build !integration

package main

import (
	"net/http"
	"testing"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
)

// nolint // need to wait for topic to be created so retry logic is to be added
func TestServerRun(t *testing.T) {
	go main()
	time.Sleep(3 * time.Second)

	tests := []struct {
		desc       string
		id         int
		method     string
		endpoint   string
		statusCode int
	}{
		{"publish", 1, http.MethodGet, "http://localhost:9113/pub?id=1", http.StatusOK},
		{"subscribe", 2, http.MethodGet, "http://localhost:9113/sub", http.StatusOK},
	}

	for i, tc := range tests {
		req, _ := request.NewMock(tc.method, tc.endpoint, nil)
		c := http.Client{}

		resp, err := c.Do(req)
		if err != nil {
			t.Errorf("TEST[%v] Failed.\tHTTP request encountered Err: %v\n%s", i, err, tc.desc)
			continue
		}

		if resp.StatusCode != tc.statusCode {
			// required because eventhub is shared and there can be messages with avro or without avro
			// messages without avro would return 200 as we do json.Marshal to a map
			// messages with avro would return 206 as it would have to go through avro.Marshal
			// we can't use any avro schema as any schema can be used
			if resp.StatusCode != http.StatusPartialContent {
				t.Errorf("TEST[%v] FAILED.\tExpected %v\tGot %v\n%s", tc.id, tc.statusCode, resp.StatusCode, tc.desc)
			}
		}

		if resp.StatusCode != tc.statusCode {
			t.Errorf("TEST[%v] FAILED.\tExpected %v\tGot %v\n%s", tc.id, tc.statusCode, resp.StatusCode, tc.desc)
		}

		_ = resp.Body.Close()

	}
}
