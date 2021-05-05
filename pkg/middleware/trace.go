package middleware

import (
	"fmt"
	"net/http"

	"go.opencensus.io/trace"
)

// Trace is a middleware which starts a span and the newly added context can be propagated and used for tracing
func Trace(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Start Context and Tracing
		ctx := r.Context()
		ctx, span := trace.StartSpan(ctx, fmt.Sprintf("%s %s", r.Method, r.URL.Path))
		defer span.End()
		inner.ServeHTTP(w, r.WithContext(ctx))
	})
}
