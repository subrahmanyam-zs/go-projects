package gofr

import (
	"context"
	"strings"

	"developer.zopsmart.com/go/gofr/pkg/log"

	"go.opentelemetry.io/collector/translator/conventions"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	cloudtrace "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
)

type exporter struct {
	name    string
	host    string
	port    string
	url     string
	appName string
}

func TraceProvider(appName, exporterName, exporterHost, exporterPort string, logger log.Logger, config Config) *trace.TracerProvider {
	exporterName = strings.ToLower(exporterName)
	e := exporter{
		name:    exporterName,
		url:     config.Get("TRACER_URL"),
		appName: appName,
	}

	switch exporterName {
	case "zipkin":
		return e.getZipkinExporter(config, logger)
	case "gcp":
		return getGCPExporter(config, exporterHost, logger)
	case "stdout":
		return stdOutTrace(config, logger)
	default:
		return nil
	}
}

func (e *exporter) getZipkinExporter(c Config, logger log.Logger) *trace.TracerProvider {
	url := e.url + "/api/v2/spans"

	exporter, err := zipkin.New(url)
	if err != nil {
		logger.Errorf("failed to initialize zipkinExporter export pipeline: %v", err)
		return nil
	}

	batcher := trace.NewBatchSpanProcessor(exporter)

	attributes := []attribute.KeyValue{
		attribute.String(conventions.AttributeTelemetrySDKLanguage, "go"),
		attribute.String(conventions.AttributeTelemetrySDKVersion, c.GetOrDefault("APP_VERSION", "Dev")),
		attribute.String(conventions.AttributeServiceName, c.GetOrDefault("APP_NAME", "Gofr-App")),
	}

	r, err := resource.New(context.Background(), resource.WithAttributes(attributes...))
	if err != nil {
		logger.Errorf("error in creating the resource")
		return nil
	}

	tp := trace.NewTracerProvider(trace.WithSampler(trace.AlwaysSample()), trace.WithSpanProcessor(batcher), trace.WithResource(r))

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp
}

func Resource() *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("stdout-example"),
		semconv.ServiceVersionKey.String("0.0.1"),
	)
}

func stdOutTrace(c Config, logger log.Logger) *trace.TracerProvider {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		logger.Errorf("creating stdout exporter: %v", err)
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(Resource()),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tracerProvider
}

func getGCPExporter(c Config, projectID string, logger log.Logger) *trace.TracerProvider {
	exporter, err := cloudtrace.New(cloudtrace.WithProjectID(projectID))
	if err != nil {
		logger.Errorf("%v", err)
		return nil
	}

	attributes := []attribute.KeyValue{
		attribute.String(conventions.AttributeTelemetrySDKLanguage, "go"),
		attribute.String(conventions.AttributeTelemetrySDKVersion, c.GetOrDefault("APP_VERSION", "Dev")),
		attribute.String(conventions.AttributeServiceName, c.GetOrDefault("APP_NAME", "Gofr-App")),
	}

	r, err := resource.New(context.Background(), resource.WithAttributes(attributes...))
	if err != nil {
		logger.Errorf("error creating resource")
		return nil
	}

	tp := trace.NewTracerProvider(
		// For this example code we use sdktrace.AlwaysSample sampler to sample all traces.
		// In a production application, use sdktrace.ProbabilitySampler with a desired probability.
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exporter),
		trace.WithResource(r))

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp
}
