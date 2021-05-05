package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/zopsmart/gofr/examples/using-pubsub/handlers"
	"github.com/zopsmart/gofr/pkg/errors"
	"github.com/zopsmart/gofr/pkg/gofr/request"
)

// nolint, need to wait for topic to be created so retry logic is to be added
func TestServerRun(t *testing.T) {
	// avro schema registry test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		re := map[string]interface{}{
			"subject": "gofr-value",
			"version": 3,
			"id":      303,
			"schema": "{\"type\":\"record\",\"name\":\"person\"," +
				"\"fields\":[{\"name\":\"Id\",\"type\":\"string\"}," +
				"{\"name\":\"Name\",\"type\":\"string\"}," +
				"{\"name\":\"Email\",\"type\":\"string\"}]}",
		}

		reBytes, _ := json.Marshal(re)
		w.Header().Set("Content-type", "application/json")
		_, _ = w.Write(reBytes)
	}))

	schemaURL := os.Getenv("AVRO_SCHEMA_URL")
	os.Setenv("AVRO_SCHEMA_URL", ts.URL)

	topic := os.Getenv("KAFKA_TOPIC")
	os.Setenv("KAFKA_TOPIC", "avro-pubsub")

	defer func() {
		os.Setenv("AVRO_SCHEMA_URL", schemaURL)
		os.Setenv("KAFKA_TOPIC", topic)
	}()

	go main()
	time.Sleep(3 * time.Second)

	tcs := []struct {
		id                 int
		method             string
		endpoint           string
		expectedResponse   string
		expectedStatusCode int
	}{
		{1, "GET", "http://localhost:9111/pub?id=1", "", 200},
		{2, "GET", "http://localhost:9111/sub", "1", 200},
	}

	for _, tc := range tcs {
		req, _ := request.NewMock(tc.method, tc.endpoint, nil)
		c := http.Client{}

		for i := 0; i < 5; i++ {
			resp, _ := c.Do(req)

			if resp != nil && resp.StatusCode != 200 {
				// retry is required since, creation of topic takes time
				if checkRetry(resp.Body) {
					time.Sleep(3 * time.Second)
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

	if len(errResp.Errors) == 0 {
		return false
	}

	_ = json.Unmarshal(body, &errResp)

	return strings.Contains(errResp.Errors[0].Reason, "Leader Not Available: the cluster is in the middle of a leadership election")
}
