package middleware

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"go.opentelemetry.io/otel/trace"
)

type MockHandlerForTracing struct{}

// ServeHTTP is used for testing if the request context has traceId
func (r *MockHandlerForTracing) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	traceID := trace.SpanFromContext(req.Context()).SpanContext().TraceID().String()
	_, _ = w.Write([]byte(traceID))
}

func TestTrace(t *testing.T) {
	handler := Trace(&MockHandlerForTracing{})
	req := httptest.NewRequest("GET", "/dummy", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, req)

	traceID := recorder.Body.String()

	if traceID == "" {
		t.Errorf("Failed to get traceId")
	}

	// if tracing has failed then the traceId is usually '00000000000000000000000000000000'
	// which is not an empty string, hence conversion to int is required to check if tracing id is correct.
	id, err := strconv.Atoi(traceID)

	if err == nil && id == 0 {
		t.Errorf("Incorrect tracingId")
	}
}