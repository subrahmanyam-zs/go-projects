package service

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_checkHealth(t *testing.T) {
	tcs := []struct {
		ts       *httptest.Server
		expected bool
	}{
		{
			httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})), false,
		},
		{
			httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})), true,
		},
	}

	for _, tc := range tcs {
		sp := new(surgeProtector)
		ch := make(chan bool)

		sp.customHeartbeatURL = "/.well-known/heartbeat"
		sp.retryFrequencySeconds = 1

		go sp.checkHealth(tc.ts.URL, ch)

		time.Sleep(5 * time.Second)

		if got := <-ch; got != tc.expected {
			t.Errorf("FAILED, Expected: %v, Got: %v", tc.expected, got)
		}

		tc.ts.Close()
	}
}
