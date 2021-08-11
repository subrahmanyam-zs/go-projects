package middleware

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

// Trace is a middleware which starts a span and the newly added context can be propagated and used for tracing
func Trace(appName,appVersion string) func(inner http.Handler) http.Handler {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			tracer := otel.GetTracerProvider().Tracer(appName, trace.WithInstrumentationVersion(appVersion))

			ctx, span := tracer.Start(ctx, fmt.Sprintf("gofr-middleware %s %s", r.Method, r.URL.Path), trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(semconv.ServiceNameKey.String("Gofr-App"), semconv.TelemetrySDKNameKey.String("Zipkin")))
			defer span.End()

			inner.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
