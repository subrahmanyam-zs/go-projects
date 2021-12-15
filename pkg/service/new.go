package service

import (
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github-lvs.corpzone.internalzone.com/mcafee/cnsr-gofr-csp-auth/generator"

	"go.opencensus.io/plugin/ochttp"

	"developer.zopsmart.com/go/gofr/pkg"
	"developer.zopsmart.com/go/gofr/pkg/gofr/cache"
	"developer.zopsmart.com/go/gofr/pkg/gofr/types"
	"developer.zopsmart.com/go/gofr/pkg/log"
)

// Options allows the user set all the options needs for http service like auth, service level headers, caching and surge protection
type Options struct {
	Headers      map[string]string // this can be used to pass service level headers.
	NumOfRetries int
	*Auth
	*Cache
	*SurgeProtectorOption
}

// KeyGenerator provides ability to  the user that can use custom or own logic to Generate the key for HTTPCached
type KeyGenerator func(url string, params map[string]interface{}, headers map[string]string) string

// Auth stores the information related to authentication. One can either use basic auth or OAuth
type Auth struct {
	UserName string // if token is not sent then the username and password can be sent and the token will be generated by the framework
	Password string
	*OAuthOption
	// Deprecated: Instead us CSPSecurityOption
	*CSPOption

	CSPSecurityOption *generator.Option
}

// Cache provides the options needed for caching of HTTPService responses
type Cache struct {
	cache.Cacher
	TTL          time.Duration
	KeyGenerator KeyGenerator
}

type SurgeProtectorOption struct {
	HeartbeatURL   string
	RetryFrequency int // indicates the time in seconds
	Disable        bool
}

//nolint // The declared global variable can be accessed across multiple functions
var (
	httpServiceResponse = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "zs_http_service_response",
		Help:    "Histogram of HTTP response times in seconds and status",
		Buckets: []float64{.001, .003, .005, .01, .025, .05, .1, .2, .3, .4, .5, .75, 1, 2, 3, 5, 10, 30},
	}, []string{"path", "method", "status"})

	circuitOpenCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "zs_external_service_circuit_open_count",
		Help: "Counter to track the number of times circuit opens",
	}, []string{"host"})
)

// nolint:gocognit,gocyclo // don't want users to access methods of this type without initialization
// hence we dont want to export the type
func NewHTTPServiceWithOptions(resourceAddr string, logger log.Logger, options *Options) *httpService {
	// Register the prometheus metric
	resourceAddr = strings.TrimRight(resourceAddr, "/")
	_ = prometheus.Register(httpServiceResponse)
	transport := &ochttp.Transport{}

	httpSvc := &httpService{
		url:       resourceAddr,
		logger:    logger,
		Client:    &http.Client{Transport: transport, Timeout: RetryFrequency * time.Second}, // default timeout is 5 seconds
		isHealthy: true,
		healthCh:  make(chan bool),
		sp: surgeProtector{
			isEnabled:             true,
			customHeartbeatURL:    "/.well-known/heartbeat",
			retryFrequencySeconds: RetryFrequency,
		},
	}

	if options == nil {
		httpSvc.SetSurgeProtectorOptions(true, httpSvc.sp.customHeartbeatURL, httpSvc.sp.retryFrequencySeconds)
		return httpSvc
	}

	// enable retries for call
	httpSvc.numOfRetries = options.NumOfRetries

	// enable service level headers
	if options.Headers != nil {
		httpSvc.customHeaders = options.Headers
	}

	// enable auth
	if options.Auth != nil && options.UserName != "" && options.OAuthOption == nil { // OAuth and basic auth cannot co-exist
		httpSvc.isAuthSet = true
		httpSvc.auth = "Basic " + base64.StdEncoding.EncodeToString([]byte(options.UserName+":"+options.Password))
	}

	// enable oauth
	if options.Auth != nil && options.OAuthOption != nil && httpSvc.auth == "" { // if auth is already set to basic auth, dont set oauth
		httpSvc.isAuthSet = true
		go httpSvc.setClientOauthHeader(options.OAuthOption)
	}

	if options.Auth != nil && options.CSPOption != nil && options.CSPSecurityOption == nil {
		logger.Warn("Deprecated CSPOption is used, instead use CSPSecurityOption for CSP Security")

		options.CSPSecurityOption = &generator.Option{
			AppKey:      options.AppKey,
			SharedKey:   options.SharedKey,
			MachineName: options.MachineName,
			IPAddress:   options.IPAddress,
		}
	}

	if options.Auth != nil && options.CSPSecurityOption != nil {
		var err error

		httpSvc.csp, err = generator.New(options.CSPSecurityOption)
		if err != nil {
			logger.Warnf("CSP Auth is not enabled, %v", err)
		}
	}

	enableSP := true

	// enable surge protection
	if options.SurgeProtectorOption != nil {
		if options.RetryFrequency != 0 {
			httpSvc.sp.retryFrequencySeconds = options.RetryFrequency
		}

		if options.HeartbeatURL != "" {
			httpSvc.sp.customHeartbeatURL = options.HeartbeatURL
		}

		enableSP = !options.SurgeProtectorOption.Disable
	}

	httpSvc.SetSurgeProtectorOptions(enableSP, httpSvc.sp.customHeartbeatURL, httpSvc.sp.retryFrequencySeconds)

	// enable http service with cache
	if options.Cache != nil {
		httpSvc.cache = &cachedHTTPService{
			httpService:  httpSvc,
			cacher:       options.Cacher,
			ttl:          options.TTL,
			keyGenerator: options.KeyGenerator,
		}
	}

	return httpSvc
}

func (h *httpService) HealthCheck() types.Health {
	h.mu.Lock()
	isHealthy := h.isHealthy
	h.mu.Unlock()

	if isHealthy {
		return types.Health{
			Name:   h.url,
			Status: pkg.StatusUp,
		}
	}

	return types.Health{
		Name:   h.url,
		Status: pkg.StatusDown,
	}
}
