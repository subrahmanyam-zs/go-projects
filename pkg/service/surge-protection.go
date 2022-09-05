package service

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/log"
)

// 1. HTTP client should not bombard the downstream service
// 2. Use a constant time period to achieve this (for ex: 5 seconds)
// 3. Let's say the HTTP client makes a request to `xyz-service/abc`, make a request to ensure the service is up
// 4. If the service is down, return 500, and asynchronously make requests after `n` seconds, where `n` is the constant
//    time period (for ex: 5 seconds)

type surgeProtector struct {
	// customHeartbeatURL is the URL that the surge protector will asynchronously call to figure out the status of a
	// service, default value is `/.well-known/heartbeat`
	customHeartbeatURL string

	// retryFrequency is the retry frequency (in seconds) of the asynchronous job that routinely checks the status
	// of services that are down
	retryFrequencySeconds int

	once sync.Once

	mu sync.Mutex

	// isEnabled is true if the surge protector is enabled
	isEnabled bool

	logger log.Logger
}

type logHealthCheck struct {
	Err        error `json:"err,omitempty"`
	StatusCode int   `json:"statusCode,omitempty"`
}

func (l logHealthCheck) String() string {
	if l.Err != nil {
		return "Health Check Failed with error: " + l.Err.Error()
	}

	return "Health Check Failed with Status Code: " + strconv.Itoa(l.StatusCode)
}

func (sp *surgeProtector) checkHealth(url string, ch chan bool) {
	for {
		var isHealthy bool

		sp.mu.Lock()

		// this will be used to log an error when the health-check/heartbeat fails
		var l logHealthCheck

		resp, err := http.Get(url + sp.customHeartbeatURL)
		if err != nil {
			l.Err = err
		} else if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
			l.StatusCode = resp.StatusCode
			resp.Body.Close()
		} else {
			isHealthy = true
			resp.Body.Close()
		}

		if !isHealthy {
			sp.logger.Errorf("%v", l)
		}

		retryFrequency := sp.retryFrequencySeconds

		sp.mu.Unlock()

		ch <- isHealthy

		time.Sleep(time.Duration(retryFrequency) * time.Second)
	}
}
