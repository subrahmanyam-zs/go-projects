package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"developer.zopsmart.com/go/gofr/examples/using-pubsub/handlers"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
)

// nolint, need to wait for topic to be created so retry logic is to be added
func TestServerRun(t *testing.T) {
	topic := os.Getenv("KAFKA_TOPIC")
	os.Setenv("KAFKA_TOPIC", "kafka-pubsub")

	defer os.Setenv("KAFKA_TOPIC", topic)

	go main()
	time.Sleep(3 * time.Second)

	tcs := []struct {
		id                 int
		method             string
		endpoint           string
		expectedResponse   string
		expectedStatusCode int
	}{
		{1, "GET", "http://localhost:9112/pub?id=1", "", 200},
		{2, "GET", "http://localhost:9112/sub", "1", 200},
	}

	for _, tc := range tcs {
		req, _ := request.NewMock(tc.method, tc.endpoint, nil)
		c := http.Client{}

		for i := 0; i < 5; i++ {
			resp, _ := c.Do(req)

			if resp != nil && resp.StatusCode != 200 {
				// retry is required since, creation of topic takes time
				if checkRetry(resp.Body) {
					time.Sleep(5 * time.Second)
					continue
				}

				t.Errorf("Test %v: Failed.\tExpected %v\tGot %v\n", tc.id, 200, resp.StatusCode)

				return
			}

			if resp != nil && resp.StatusCode != tc.expectedStatusCode {
				t.Errorf("Test %v: Failed.\tExpected %v\tGot %v\n", tc.id, tc.expectedStatusCode, resp.StatusCode)
			}

			// checks whether bind avro.Unmarshal functionality works fine
			if tc.expectedResponse != "" && resp.Body != nil {
				body, _ := io.ReadAll(resp.Body)

				m := struct {
					Data handlers.Person `json:"data"`
				}{}
				_ = json.Unmarshal(body, &m)

				if m.Data.ID != tc.expectedResponse {
					t.Errorf("Expected: %v, Got: %v", tc.expectedResponse, m.Data.ID)
				}
			}

			resp.Body.Close()

			break
		}
	}
}

func checkRetry(respBody io.Reader) bool {
	body, _ := io.ReadAll(respBody)

	errResp := struct {
		Errors []errors.Response `json:"errors"`
	}{}

	_ = json.Unmarshal(body, &errResp)

	return strings.Contains(errResp.Errors[0].Reason, "Leader Not Available: the cluster is in the middle of a leadership election")
}
