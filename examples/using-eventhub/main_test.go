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

	tcs := []struct {
		id                 int
		method             string
		endpoint           string
		expectedStatusCode int
	}{
		{1, "GET", "http://localhost:9113/pub?id=1", 200},
		{2, "GET", "http://localhost:9113/sub", 200},
	}

	for _, tc := range tcs {
		req, _ := request.NewMock(tc.method, tc.endpoint, nil)
		c := http.Client{}

		resp, _ := c.Do(req)

		if resp != nil && resp.StatusCode != 200 {
			// required because eventhub is shared and there can be messages with avro or without avro
			// messages without avro would return 200 as we do json.Marshal to a map
			// messages with avro would return 206 as it would have to go through avro.Marshal
			// we can't use any avro schema as any schema can be used
			if resp.StatusCode != 206 {
				t.Errorf("Test %v: Failed.\tExpected %v\tGot %v\n", tc.id, tc.expectedStatusCode, resp.StatusCode)
			}
		}

		if resp != nil && resp.StatusCode != tc.expectedStatusCode {
			t.Errorf("Test %v: Failed.\tExpected %v\tGot %v\n", tc.id, tc.expectedStatusCode, resp.StatusCode)
		}

		if resp != nil {
			resp.Body.Close()
		}

		break
	}
}
