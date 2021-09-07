package main

import (
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
)

func TestIntegration(t *testing.T) {
	go main()
	time.Sleep(3 * time.Second)

	tcs := []struct {
		endpoint           string
		expectedStatusCode int
	}{
		{"user/ ", http.StatusBadRequest},
		{"dummyendpoint", http.StatusNotFound},
	}

	for _, tc := range tcs {
		req, _ := request.NewMock(http.MethodGet, "http://localhost:9091/"+tc.endpoint, nil)
		c := http.Client{}

		resp, _ := c.Do(req)

		if resp != nil {
			if resp.StatusCode != tc.expectedStatusCode {
				t.Errorf("Failed.\tExpected %v\tGot %v\n", tc.expectedStatusCode, resp.StatusCode)
			}

			bodyBytes, _ := ioutil.ReadAll(resp.Body)

			if reflect.DeepEqual(bodyBytes, nil) {
				t.Errorf("Failed.\tExpected %v\tGot %v\n", tc.expectedStatusCode, resp.StatusCode)
			}
		}

		resp.Body.Close()
	}
}
