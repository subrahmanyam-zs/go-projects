package middleware

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/gorilla/mux"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/log"
)

type MockHandler struct {
	statusCode int
}
type MockWriteHandler struct {
}

func (m MockWriteHandler) Header() http.Header {
	return http.Header{}
}
func (m MockWriteHandler) Write(b []byte) (int, error) {
	return 0, nil
}
func (m MockWriteHandler) WriteHeader(statuscode int) {}

func (r *MockHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if r.statusCode == 0 {
		r.statusCode = http.StatusOK
	}

	w.WriteHeader(r.statusCode)
	_, _ = w.Write([]byte("testing log"))
}

func TestLogging(t *testing.T) {
	b := new(bytes.Buffer)
	logger := log.NewMockLogger(b)
	handler := Logging(logger, "")(&MockHandler{})

	req := httptest.NewRequest("GET", "/dummy", nil)
	handler.ServeHTTP(MockWriteHandler{}, req)

	if len(b.Bytes()) == 0 {
		t.Errorf("Failed to write the logs")
	}

	x := b.String()
	if !strings.Contains(x, "time") || !strings.Contains(x, "level") {
		t.Errorf("error, expected fields are not present in log, got: %v", x)
	}
}

func TestGetIPAddress(t *testing.T) {
	{
		// 1. When RemoteAddr is set
		addr := "0.0.0.0:8080"
		req, err := http.NewRequest("GET", "http://dummy", nil)
		if err != nil {
			t.Errorf("FAILED, got error creating req object: %v", err)
		}

		req.RemoteAddr = addr

		if ip := GetIPAddress(req); ip != addr {
			t.Errorf("FAILED, expected: %v, got: %v", addr, ip)
		}
	}

	{
		// 2. When `X-Forwarded-For` header is set
		addr := "192.168.0.1:8080"
		req, err := http.NewRequest("GET", "http://dummy", nil)
		if err != nil {
			t.Errorf("FAILED, got error creating req object: %v", err)
		}

		req.Header.Set("X-Forwarded-For", addr)

		if ip := GetIPAddress(req); ip != addr {
			t.Errorf("FAILED, expected: %v, got: %v", addr, ip)
		}
	}
}

func TestLoggingCorrelationID(t *testing.T) {
	b := new(bytes.Buffer)
	logger := log.NewMockLogger(b)

	handler := Logging(logger, "")(&MockHandler{})

	req := httptest.NewRequest("GET", "/dummy", nil)
	req.Header.Add("X-B3-TraceId", "12bhu987")
	handler.ServeHTTP(MockWriteHandler{}, req)

	if len(b.Bytes()) == 0 {
		t.Errorf("Failed to write the logs")
	}

	x := b.String()
	if !strings.Contains(x, "correlationId") || !strings.Contains(x, "12bhu987") {
		t.Errorf("error, expected correlation id in log, got: %v", x)
	}
}

func TestLoggingCorrelationContext(t *testing.T) {
	b := new(bytes.Buffer)
	logger := log.NewMockLogger(b)

	handler := Logging(logger, "")(&MockHandler{})

	correlationID := "12bhu987"

	req := httptest.NewRequest("GET", "/dummy", nil)
	req.Header.Add("X-Correlation-Id", correlationID)
	handler.ServeHTTP(MockWriteHandler{}, req)

	cID, _ := req.Context().Value(CorrelationIDKey).(string)

	if cID != correlationID {
		t.Errorf("correlationID is not present in the request context")
	}
}

func TestLoggingOmitHeader(t *testing.T) {
	b := new(bytes.Buffer)
	logger := log.NewMockLogger(b)
	omitHeaders := "X-Some-Random-Header-1,X-Some-Random-Header-2,X-Some-Random-Header-3"
	handler := Logging(logger, omitHeaders)(&MockHandler{})

	req := httptest.NewRequest("GET", "/dummy", nil)
	req.Header.Add("X-Some-Random-Header-1", "Some-Random-Value")
	req.Header.Add("X-Some-Random-Header-2", "Some-Random-Value")
	req.Header.Add("X-Some-random-header-3", "Some-Random-Value-Case-Insensitive")
	handler.ServeHTTP(MockWriteHandler{}, req)

	if len(b.Bytes()) == 0 {
		t.Errorf("Failed to write the logs")
	}

	x := b.String()
	if !strings.Contains(x, "X-Some-Random-Header-1") || !strings.Contains(x, "X-Some-Random-Header-2") ||
		strings.Contains(x, "Some-Random-Value") || !strings.Contains(x, "xxx-masked-value-xxx") {
		t.Errorf("error, expected X-Some-Random-Header-1 and X-Some-Random-Header-1 with value :"+
			" xxx-masked-value-xxx, got: %v", x)
	}
}

func TestLoggingAuthorizationHeader(t *testing.T) {
	b := new(bytes.Buffer)
	logger := log.NewMockLogger(b)
	handler := Logging(logger, "")(&MockHandler{})

	req := httptest.NewRequest("GET", "/dummy", nil)
	req.Header.Add("Authorization", "Basic dXNlcjpwYXNz")
	handler.ServeHTTP(MockWriteHandler{}, req)

	if len(b.Bytes()) == 0 {
		t.Errorf("Failed to write the logs")
	}

	// Authorization header should be present
	x := b.String()

	if !strings.Contains(b.String(), "Authorization") || !strings.Contains(b.String(), "user") {
		t.Errorf("error, expected Authorization:user in header, got: %v", x)
	}

	// Authorization header should not be present as the auth token is invalid
	b.Reset()

	req = httptest.NewRequest("GET", "/dummy", nil)
	req.Header.Add("Authorization", "dummy")
	handler.ServeHTTP(MockWriteHandler{}, req)

	if len(b.Bytes()) == 0 {
		t.Errorf("Failed to write the logs")
	}

	x = b.String()

	if strings.Contains(x, "Authorization") {
		t.Errorf("error, Authorization Header should not be present in logs, got: %v", x)
	}

	// Authorization header should be masked
	b.Reset()

	handler = Logging(logger, "Authorization")(&MockHandler{})

	req = httptest.NewRequest("GET", "/dummy", nil)
	req.Header.Add("Authorization", "dummy")
	handler.ServeHTTP(MockWriteHandler{}, req)

	if len(b.Bytes()) == 0 {
		t.Errorf("Failed to write the logs")
	}

	x = b.String()

	if !strings.Contains(x, `"Authorization":"xxx-masked-value-xxx"`) {
		t.Errorf("error, Authorization Header should be masked in logs, got: %v", x)
	}
}

func TestAppData(t *testing.T) {
	b := new(bytes.Buffer)
	logger := log.NewMockLogger(b)

	handler := Logging(logger, "")(&MockHandler{})

	var appData LogDataKey = "appLogData"

	{
		data := &sync.Map{}
		data.Store("key1", "val1")
		req := httptest.NewRequest("GET", "/dummy", nil)
		req = req.WithContext(context.WithValue(req.Context(), appData, data))

		handler.ServeHTTP(MockWriteHandler{}, req)

		if len(b.Bytes()) == 0 {
			t.Errorf("Failed to write the logs")
		}

		x := b.String()

		if !strings.Contains(b.String(), `"data":{"key1":"val1"}}`) {
			t.Errorf("error, expected \"data\":{\"key1\":\"val1\"},\n got: %v", x)
		}
	}

	{
		b.Reset()
		data := &sync.Map{}
		data.Store("key2", "val2")
		req := httptest.NewRequest("GET", "/dummy", nil)
		req = req.WithContext(context.WithValue(req.Context(), appData, data))

		handler.ServeHTTP(MockWriteHandler{}, req)

		if len(b.Bytes()) == 0 {
			t.Errorf("Failed to write the logs")
		}

		x := b.String()

		if !strings.Contains(b.String(), `"data":{"key2":"val2"}}`) {
			t.Errorf("error, expected \"data\":{\"key2\":\"val2\"}, got: %v", x)
		}
	}
}

func Test_getUsernameForBasicAuth(t *testing.T) {
	type args struct {
		authHeader string
	}

	tests := []struct {
		name     string
		args     args
		wantUser string
		wantPass string
		wantErr  bool
	}{
		{"success", args{authHeader: "Basic dXNlcjpwYXNz"}, "user", "pass", false},
		{"invalid token", args{authHeader: "Basic a"}, "", "", true},
		{"failure", args{authHeader: "fail"}, "", "", true},
	}

	for _, tt := range tests {
		gotUser := getUsernameForBasicAuth(tt.args.authHeader)

		if gotUser != tt.wantUser {
			t.Errorf("getUsernameForBasicAuth() got = %v, want %v", gotUser, tt.wantUser)
		}
	}
}

// Test_ValidAppDataInConcurrentRequest tries to mimic the behavior of ApacheBench(ab)
// test with parameter n=15, c=5
func Test_ValidAppDataInConcurrentRequest(t *testing.T) {
	conReq := 5
	totalReq := 15
	b := new(Buffer)
	logger := log.NewMockLogger(b)
	handler := Logging(logger, "")(&MockHandlerLogging{})
	muxRouter := mux.NewRouter()

	muxRouter.NewRoute().Path("/hello-planet").Methods("GET").Handler(handler)
	muxRouter.NewRoute().Path("/hello-galaxy").Methods("GET").Handler(handler)

	var wg sync.WaitGroup

	batch := totalReq / conReq
	for i := 0; i < batch; i++ {
		wg.Add(1)

		go makeRequestPlanet(t, handler, &wg, "/hello-planet", conReq)
		wg.Add(1)

		go makeRequestGalaxy(t, handler, &wg, "/hello-galaxy", conReq)
		wg.Wait()
	}

	checkLogs(t, b)
}

func TestErrorMessages(t *testing.T) {
	b := new(bytes.Buffer)
	logger := log.NewMockLogger(b)

	errorMessage := "test-error"

	err := errors.Response{Reason: errorMessage}

	req := httptest.NewRequest("GET", "/dummy", nil)
	req = req.WithContext(context.WithValue(req.Context(), ErrorMessage, err.Error()))

	handler := Logging(logger, "")(&MockHandler{statusCode: http.StatusInternalServerError})

	handler.ServeHTTP(MockWriteHandler{}, req)

	actual := b.String()

	if !strings.Contains(actual, errorMessage) {
		t.Errorf("FAILED, expected: %v, got: %v", errorMessage, b.String())
	}
}

// TestCookieLogging checks Cookie is getting logged or not.
func TestCookieLogging(t *testing.T) {
	b := new(bytes.Buffer)
	logger := log.NewMockLogger(b)

	handler := Logging(logger, "")(&MockHandler{})

	req := httptest.NewRequest("GET", "http://dummy", nil)
	req.Header.Add("Cookie", "Some-Random-Value")

	handler.ServeHTTP(MockWriteHandler{}, req)

	x := b.String()
	if strings.Contains(x, "Cookie") {
		t.Errorf("Error: Expected no cookie, Got: %v", x)
	}
}
