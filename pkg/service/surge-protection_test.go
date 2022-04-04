package service

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/log"
)

func Test_checkHealth(t *testing.T) {
	tcs := []struct {
		statusCode      int
		expected        bool
		expectedLogData string
	}{
		{http.StatusInternalServerError, false, "Health Check Failed with Status Code: 500"},
		{http.StatusOK, true, ""},
	}

	for _, tc := range tcs {
		b := new(bytes.Buffer)
		logger := log.NewMockLogger(b)

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(tc.statusCode)
		}))

		sp := new(surgeProtector)
		ch := make(chan bool)

		sp.customHeartbeatURL = "/.well-known/heartbeat"
		sp.logger = logger

		go sp.checkHealth(ts.URL, ch)

		time.Sleep(5 * time.Second)

		if got := <-ch; got != tc.expected {
			t.Errorf("FAILED, Expected: %v, Got: %v", tc.expected, got)
		}

		if !strings.Contains(b.String(), tc.expectedLogData) {
			t.Errorf("FAILED expected %v,got: %v", tc.expectedLogData, b.String())
		}

		ts.Close()
	}
}
