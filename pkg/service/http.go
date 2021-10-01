package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.opencensus.io/plugin/ochttp"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/log"
	"developer.zopsmart.com/go/gofr/pkg/middleware"
)

type httpService struct {
	*http.Client

	auth          string
	url           string
	logger        log.Logger
	headerKeys    []string
	customHeaders map[string]string
	healthCh      chan bool
	isHealthy     bool
	isAuthSet     bool
	contentType   responseType
	sp            surgeProtector
	numOfRetries  int

	cache *cachedHTTPService
}

type responseType int

const (
	JSON responseType = iota
	XML
	TEXT
	HTML
	RetryFrequency          = 5
	UnixTimeStampMultiplier = 1000
	ErrToken                = errors.Error("token could not be obtained")
)

// nolint:lll,gocognit,gocyclo  // cannot reduce the number of lines since there are many parameters.
func (h *httpService) call(ctx context.Context, method, target string, params map[string]interface{}, body []byte, headers map[string]string) (*Response, error) {
	correlationID, _ := ctx.Value(middleware.CorrelationIDKey).(string)
	appData := getAppData(ctx)

	// store the authorization header for Auth name
	var authorizationHeader string
	if val := ctx.Value(middleware.AuthorizationHeader); val != nil {
		authorizationHeader = val.(string)
	}

	err := h.preCall(method, target, correlationID, params, appData, headers, authorizationHeader)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, RequestCanceled{}
	default:
		req, err := h.createReq(ctx, method, target, params, body, headers)
		if err != nil {
			return nil, err
		}

		headers := make(map[string]string)

		for head := range req.Header {
			val := req.Header.Get(head)
			if val != "" {
				headers[head] = req.Header.Get(head)
			}
		}

		// Don't want to log the Cookie.
		delete(headers, "Cookie")

		start := time.Now()

		var (
			resp       *http.Response
			statusCode int
		)

		for i := 0; i <= h.numOfRetries; i++ {
			resp, err = h.Do(req.WithContext(ctx)) //nolint:bodyclose // body is being closed after call response is logged

			if resp != nil {
				statusCode = resp.StatusCode
			}

			if err != nil {
				h.logError(&errorLog{CorrelationID: correlationID, Method: method, URI: h.url + "/" + target,
					ResponseCode: statusCode, Params: params, Message: err.Error(), AppData: appData}, headers, start, authorizationHeader)

				if e, ok := err.(net.Error); ok && e.Timeout() {
					// the error occurred due to timeout, so continue
					continue
				}
			}

			break
		}
		// add url, method, statusCode and duration in prometheus metric
		httpServiceResponse.WithLabelValues(h.url+"/"+target, method, fmt.Sprintf("%d", statusCode)).Observe(time.Since(start).Seconds())

		if err != nil {
			return nil, err
		}

		switch resp.Header.Get("content-type") {
		case "application/xml":
			h.contentType = XML
		case "text/plain":
			h.contentType = TEXT
		default:
			h.contentType = JSON
		}

		h.logCall(&callLog{CorrelationID: correlationID, Method: method, URI: h.url + "/" + target,
			ResponseCode: resp.StatusCode, Params: params, AppData: appData}, headers, start, authorizationHeader)

		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()

		response := Response{
			Body:       body,
			StatusCode: resp.StatusCode,
			headers:    resp.Header,
		}

		return &response, err
	}
}

// nolint:lll //cannot reduce the number of lines since there are many parameters.
func (h *httpService) preCall(method, target, correlationID string, params, appData map[string]interface{}, headers map[string]string, authorizationHeader string) error {
	start := time.Now()

	var (
		err        error
		statusCode int
	)

	if h.isAuthSet && h.auth == "" {
		err = ErrToken
		statusCode = http.StatusUnauthorized
	}

	if !h.isHealthy {
		err = ErrServiceDown{URL: h.url}
		statusCode = http.StatusInternalServerError
	}

	if err != nil {
		httpServiceResponse.WithLabelValues(h.url+"/"+target, method, fmt.Sprintf("%d", statusCode)).Observe(time.Since(start).Seconds())
		h.logError(&errorLog{CorrelationID: correlationID, Method: method, URI: h.url + "/" + target, Params: params,
			Message: err.Error(), AppData: appData}, headers, start, authorizationHeader)

		return err
	}

	return nil
}

// fetch the appData from request context and generate a map of type map[string]interface{}, if appData is nil
// then getAppData will return empty map
func getAppData(ctx context.Context) map[string]interface{} {
	appData := make(map[string]interface{})

	if data, ok := ctx.Value(middleware.LogDataKey("appLogData")).(*sync.Map); ok {
		data.Range(func(key, value interface{}) bool {
			if k, ok := key.(string); ok {
				appData[k] = value
			}

			return true
		})
	}

	return appData
}

// createReq creates the request for service call
// the endpoint and the method for the request are defined from the parameters provided to the function
// nolint:lll,gocognit // cannot reduce the number of lines since there are many parameters.
func (h *httpService) createReq(ctx context.Context, method, target string, params map[string]interface{}, body []byte, headers map[string]string) (*http.Request, error) {
	target = strings.TrimLeft(target, "/")
	uri := h.url + "/" + target

	if target == "" {
		uri = h.url
	}

	req, err := http.NewRequest(method, uri, bytes.NewBuffer(body))
	if err != nil {
		return nil, FailedRequest{URL: h.url, Err: err}
	}

	setContentTypeAndAcceptHeader(req, body)

	h.setHeadersFromContext(ctx, req)

	// service level headers
	for k, v := range h.customHeaders {
		req.Header.Set(k, v)
	}

	// request level headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// query parameters is required for GET,POST and PUT method
	if (method == "GET" || method == "POST" || method == "PUT" || method == "PATCH") && params != nil {
		encodeQueryParameters(req, params)
	}

	return req, nil
}

func (h *httpService) setHeadersFromContext(ctx context.Context, req *http.Request) {
	// add all the mandatory headers to the request
	if val := ctx.Value(middleware.CorrelationIDKey); val != nil {
		correlationID, _ := val.(string)
		req.Header.Add("X-Correlation-Id", correlationID)
	}

	if val := ctx.Value(middleware.ClientIPKey); val != nil {
		clientIP, _ := val.(string)
		req.Header.Add("True-Client-IP", clientIP)
	}

	if val := ctx.Value(middleware.ZopsmartChannelKey); val != nil {
		zopsmartChannel, _ := val.(string)
		req.Header.Add("X-Zopsmart-Channel", zopsmartChannel)
	}

	if val := ctx.Value(middleware.AuthenticatedUserIDKey); val != nil {
		authUserID, _ := val.(string)
		req.Header.Add("X-Authenticated-UserId", authUserID)
	}

	if val := ctx.Value(middleware.ZopsmartTenantKey); val != nil {
		zopsmartTenant, _ := val.(string)
		req.Header.Add("X-Zopsmart-Tenant", zopsmartTenant)
	}

	if h.auth != "" {
		req.Header.Add("Authorization", h.auth)
	}

	// add custom headers to the request
	for i := range h.headerKeys {
		val, _ := ctx.Value(h.headerKeys[i]).(string)
		req.Header.Add(h.headerKeys[i], val)
	}
}

func setContentTypeAndAcceptHeader(req *http.Request, body []byte) {
	if body == nil {
		req.Header.Set("Accept", "application/json,application/xml,text/plain")
		return
	}

	contentType := "text/plain"

	var temp interface{}

	err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&temp)
	if err == nil {
		contentType = "application/json"
	}

	err = xml.NewDecoder(bytes.NewBuffer(body)).Decode(&temp)
	if err == nil {
		contentType = "application/xml"
	}

	req.Header.Add("content-type", contentType)
}

func encodeQueryParameters(req *http.Request, params map[string]interface{}) {
	q := req.URL.Query()

	for k, v := range params {
		switch vt := v.(type) {
		case []string:
			for _, val := range vt {
				q.Set(k, val)
			}
		default:
			q.Set(k, fmt.Sprintf("%v", v))
		}
	}

	req.URL.RawQuery = q.Encode()
}

func (h *httpService) SetConnectionPool(maxConnections int, idleConnectionTimeout time.Duration) {
	t := http.Transport{MaxIdleConns: maxConnections, IdleConnTimeout: idleConnectionTimeout}
	octr := &ochttp.Transport{Base: &t}
	h.Timeout = idleConnectionTimeout
	cl := &http.Client{Transport: octr}
	h.Client = cl
}

func (h *httpService) PropagateHeaders(headers ...string) {
	for i := range headers {
		h.headerKeys = append(h.headerKeys, headers[i])
	}
}

// nolint:gocognit // Splitting the code will reduce readability
func (h *httpService) SetSurgeProtectorOptions(isEnabled bool, customHeartbeatURL string, retryFrequencySeconds int) {
	if isEnabled {
		// Register the prometheus metric
		_ = prometheus.Register(circuitOpenCount)

		h.sp.once.Do(func() {
			go func() {
				go h.sp.checkHealth(h.url, h.healthCh)

				for ok := range h.healthCh {
					h.sp.mu.Lock()
					// If the circuit is open, the circuitOpenCount metric value will be increased otherwise, value will not change
					if !ok && h.isHealthy {
						circuitOpenCount.WithLabelValues(h.url).Inc()
					}

					h.isHealthy = ok
					h.sp.mu.Unlock()
				}
			}()
		})
	}

	h.sp.mu.Lock()
	defer h.sp.mu.Unlock()

	h.sp.isEnabled = isEnabled

	if customHeartbeatURL != "" {
		h.sp.customHeartbeatURL = customHeartbeatURL
	}

	if retryFrequencySeconds != 0 {
		h.sp.retryFrequencySeconds = retryFrequencySeconds
	}
}

type callLog struct {
	CorrelationID string                 `json:"correlationId"`
	Type          string                 `json:"type"`
	Timestamp     int64                  `json:"timestamp"`
	Duration      int64                  `json:"duration"`
	Method        string                 `json:"method"`
	URI           string                 `json:"uri"`
	ResponseCode  int                    `json:"responseCode"`
	Params        map[string]interface{} `json:"params,omitempty"`
	Headers       map[string]string      `json:"headers,omitempty"`
	AppData       map[string]interface{} `json:"appData,omitempty"`
}

type errorLog struct {
	CorrelationID string                 `json:"correlationId"`
	Type          string                 `json:"type"`
	Timestamp     int64                  `json:"timestamp"`
	Duration      int64                  `json:"duration"`
	Method        string                 `json:"method"`
	URI           string                 `json:"uri"`
	ResponseCode  int                    `json:"responseCode,omitempty"`
	Params        map[string]interface{} `json:"params,omitempty"`
	Headers       map[string]string      `json:"headers,omitempty"`
	Message       string                 `json:"message"`
	AppData       map[string]interface{} `json:"appData,omitempty"`
}

func (l *callLog) String() string {
	line, _ := json.Marshal(l)
	return string(line)
}

func (l *errorLog) String() string {
	line, _ := json.Marshal(l)
	return string(line)
}

func (h *httpService) logCall(l *callLog, headers map[string]string, startTime time.Time, authorizationHeader string) {
	setAuthHeader(headers, authorizationHeader)
	l.Headers = headers
	l.Type = "PERFORMANCE"
	l.Duration = time.Since(startTime).Microseconds()
	l.Timestamp = startTime.Unix() * UnixTimeStampMultiplier

	if h.logger != nil {
		h.logger.Logf("%v", l)
	}
}

func (h *httpService) logError(l *errorLog, headers map[string]string, startTime time.Time, authorizationHeader string) {
	setAuthHeader(headers, authorizationHeader)
	l.Headers = headers
	l.Type = "ERROR"
	l.Duration = time.Since(startTime).Microseconds()
	l.Timestamp = startTime.Unix() * UnixTimeStampMultiplier

	if h.logger != nil {
		h.logger.Errorf("%v", l)
	}
}

func setAuthHeader(headers map[string]string, authorizationHeader string) {
	delete(headers, "Authorization")

	name, err := getUsername(authorizationHeader)
	if err == nil {
		headers["Authorization"] = name
	}
}

func getUsername(authHeader string) (user string, err error) {
	const authLen = 2
	auth := strings.SplitN(authHeader, " ", authLen)

	if authHeader == "" {
		return "", middleware.ErrMissingHeader
	}

	if len(auth) != authLen || auth[0] != "Basic" {
		return "", middleware.ErrInvalidHeader
	}

	payload, _ := base64.StdEncoding.DecodeString(auth[1])
	pair := strings.SplitN(string(payload), ":", authLen)

	if len(pair) < authLen {
		return "", middleware.ErrInvalidToken
	}

	return pair[0], nil
}
