package middleware

import (
	"fmt"
	"net/http"

	"developer.zopsmart.com/go/gofr/pkg"

	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

// Trace is a middleware which starts a span and the newly added context can be propagated and used for tracing
func Trace(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		tracer := otel.GetTracerProvider().Tracer(pkg.DefaultAppName, trace.WithInstrumentationVersion(pkg.DefaultAppVersion))

		ctx, span := tracer.Start(ctx, fmt.Sprintf("gofr-middleware %s %s", r.Method, r.URL.Path), trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(semconv.ServiceNameKey.String("Gofr-App"), semconv.TelemetrySDKNameKey.String("Zipkin")))
		defer span.End()

		inner.ServeHTTP(w, r.WithContext(ctx))
	})
}
