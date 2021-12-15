package gofr

import (
	"context"
	"strings"

	cloudtrace "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"

	"go.opentelemetry.io/collector/translator/conventions"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"

	"developer.zopsmart.com/go/gofr/pkg/errors"
)

type exporter struct {
	name    string
	url     string
	appName string
}

func tracerProvider(config Config) (err error) {
	appName := config.GetOrDefault("APP_NAME", "gofr")
	exporterName := strings.ToLower(config.Get("TRACER_EXPORTER"))
	exporterHost := config.Get("TRACER_HOST")

	e := exporter{
		name:    exporterName,
		url:     config.Get("TRACER_URL"),
		appName: appName,
	}

	var tp *trace.TracerProvider

	switch exporterName {
	case "zipkin":
		tp, err = e.getZipkinExporter(config)
	case "gcp":
		tp, err = getGCPExporter(config, exporterHost)
	case "stdout":
		tp, err = stdOutTrace(config)
	default:
		return errors.Error("invalid exporter")
	}

	if err != nil {
		return
	}

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return
}

func (e *exporter) getZipkinExporter(c Config) (*trace.TracerProvider, error) {
	url := e.url + "/api/v2/spans"

	exporter, err := zipkin.New(url)
	if err != nil {
		return nil, err
	}

	batcher := trace.NewBatchSpanProcessor(exporter)

	r, err := getResource(c)
	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(trace.WithSampler(trace.AlwaysSample()), trace.WithSpanProcessor(batcher), trace.WithResource(r))

	return tp, nil
}

func stdOutTrace(c Config) (*trace.TracerProvider, error) {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	r, err := getResource(c)
	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(trace.WithBatcher(exporter), trace.WithResource(r))

	return tp, nil
}

func getGCPExporter(c Config, projectID string) (*trace.TracerProvider, error) {
	exporter, err := cloudtrace.New(cloudtrace.WithProjectID(projectID))
	if err != nil {
		return nil, err
	}

	r, err := getResource(c)
	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exporter),
		trace.WithResource(r))

	return tp, nil
}

func getResource(c Config) (*resource.Resource, error) {
	attributes := []attribute.KeyValue{
		attribute.String(conventions.AttributeTelemetrySDKLanguage, "go"),
		attribute.String(conventions.AttributeTelemetrySDKVersion, c.GetOrDefault("APP_VERSION", "Dev")),
		attribute.String(conventions.AttributeServiceName, c.GetOrDefault("APP_NAME", "Gofr-App")),
	}

	return resource.New(context.Background(), resource.WithAttributes(attributes...))
}
