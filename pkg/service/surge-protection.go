package service

import (
	"net/http"
	"sync"
	"time"
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
}

func (sp *surgeProtector) checkHealth(url string, ch chan bool) {
	for {
		var isHealthy bool

		resp, err := http.Get(url + sp.customHeartbeatURL)
		if err != nil || (resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound) {
			isHealthy = false
		} else {
			isHealthy = true
			resp.Body.Close()
		}

		ch <- isHealthy

		time.Sleep(time.Duration(sp.retryFrequencySeconds) * time.Second)
	}
}
