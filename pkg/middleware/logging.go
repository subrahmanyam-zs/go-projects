package middleware

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"go.opentelemetry.io/otel/trace"
)

type message string

const ErrorMessage message = "errorMessage"

// StatusResponseWriter defines own Response Writer to be used for logging of status - as http.ResponseWriter does not let us read status.
type StatusResponseWriter struct {
	http.ResponseWriter
	status int
}

type LogDataKey string

func (w *StatusResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

type LogLine struct {
	CorrelationID  string                 `json:"correlationId"`
	Type           string                 `json:"type"`
	StartTimestamp time.Time              `json:"startTimestamp"`
	Duration       int64                  `json:"duration"`
	Method         string                 `json:"method"`
	IP             string                 `json:"ip"`
	URI            string                 `json:"uri"`
	Response       int                    `json:"responseCode"`
	Headers        map[string]string      `json:"headers"`
	AppData        map[string]interface{} `json:"appData"`
	ErrorMessage   string                 `json:"errorMessage,omitempty"`
}

func (l *LogLine) String() string {
	line, _ := json.Marshal(l)
	return string(line)
}

type logger interface {
	Log(a ...interface{})
	Debug(args ...interface{})
	AddData(key string, value interface{})
	Errorf(format string, a ...interface{})
}

type contextKey string

const CorrelationIDKey contextKey = "correlationID"

// It's sequential statements and some closures. Trying to break it will make it
// Logging is a middleware which logs response status and time in microseconds along with other data.
func Logging(logger logger, omitHeaders string) func(inner http.Handler) http.Handler {
	omitHeadersMap := getOmitLogHeader(omitHeaders)

	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			correlationID := GetCorrelationID(r)
			ctx := context.WithValue(r.Context(), CorrelationIDKey, correlationID)
			*r = *r.WithContext(ctx)

			srw := &StatusResponseWriter{ResponseWriter: w}
			defer func(res *StatusResponseWriter, req *http.Request) {
				headers := fetchHeaders(omitHeadersMap, req.Header)

				l := LogLine{
					CorrelationID:  correlationID,
					StartTimestamp: start,
					Duration:       time.Since(start).Microseconds(),
					Method:         req.Method,
					IP:             GetIPAddress(req),
					URI:            req.RequestURI,
					Response:       res.status,
					Type:           "PERFORMANCE",
					Headers:        headers,
				}

				l.ErrorMessage = populateMessage(r, res.status)

				if logger != nil {
					// fetch the appData from request context and generate a map of type map[string]interface{}, if appData is nil
					// then getAppData will return empty map
					l.AppData = getAppData(req.Context())

					// .well-known, swagger and metrics endpoints are logged in debug mode, so that it can be excluded
					// from logs, as usually logs with level INFO or higher than INFO are logged
					if ExemptPath(r) {
						logger.Debug(&l)
					} else {
						logger.Log(&l)
					}
				}
			}(srw, r)

			inner.ServeHTTP(srw, r)
		})
	}
}

func GetCorrelationID(r *http.Request) string {
	correlationID := r.Header.Get("X-Correlation-Id")
	if correlationID == "" {
		correlationID = r.Header.Get("X-B3-TraceId")
	}

	if correlationID == "" {
		correlationID = trace.SpanFromContext(r.Context()).SpanContext().TraceID().String()
		r.Header.Set("X-Correlation-Id", correlationID)
	}

	return correlationID
}

func GetIPAddress(r *http.Request) string {
	var ipAddress string

	ips := strings.Split(r.Header.Get("X-Forwarded-For"), ",")

	// According to GCLB Documentation (https://cloud.google.com/load-balancing/docs/https/), IPs are added in following sequence.
	// X-Forwarded-For: <unverified IP(s)>, <immediate client IP>, <global forwarding rule external IP>, <proxies running in GCP>
	ipAddress = ips[0]

	if ipAddress == "" {
		ipAddress = r.RemoteAddr
	}

	return strings.TrimSpace(ipAddress)
}

func fetchHeaders(omitHeaders map[string]bool, reqHeaders http.Header) map[string]string {
	headers := make(map[string]string)

	for h := range reqHeaders {
		lowerCase := strings.ToLower(h)
		if _, ok := omitHeaders[lowerCase]; !ok {
			if lowerCase == "authorization" {
				processAuthHeader(headers, h, reqHeaders.Get(h))
			} else if lowerCase == "ac" || lowerCase == "ak" {
				// Don't want to log the CSP ac and ak headers.
				continue
			} else {
				headers[h] = reqHeaders.Get(h)
			}
		} else {
			headers[h] = "xxx-masked-value-xxx"
		}
	}

	// Don't want to log the Cookie.
	delete(headers, "Cookie")

	return headers
}

func processAuthHeader(headers map[string]string, authHeader, value string) {
	userName := getUsernameForBasicAuth(value)
	if userName != "" {
		headers[authHeader] = userName
	}
}

func getOmitLogHeader(headers string) map[string]bool {
	omitHeadersMap := make(map[string]bool)

	headersList := strings.Split(headers, ",")
	if len(headersList) == 1 && headersList[0] == "" {
		return omitHeadersMap
	}

	for _, h := range headersList {
		// for case insensitive headers
		lowerCase := strings.ToLower(h)
		omitHeadersMap[lowerCase] = true
	}

	return omitHeadersMap
}

func getUsernameForBasicAuth(authHeader string) (user string) {
	const authLen = 2
	auth := strings.SplitN(authHeader, " ", authLen)

	if authHeader == "" {
		return ""
	}

	if len(auth) != authLen || auth[0] != "Basic" {
		return ""
	}

	payload, _ := base64.StdEncoding.DecodeString(auth[1])
	pair := strings.SplitN(string(payload), ":", authLen)

	if len(pair) < authLen {
		return ""
	}

	return pair[0]
}

func populateMessage(r *http.Request, statusCode int) string {
	var msg string

	if statusCode < http.StatusOK || statusCode >= http.StatusMultipleChoices {
		msg, _ = r.Context().Value(ErrorMessage).(string)
	}

	return msg
}
