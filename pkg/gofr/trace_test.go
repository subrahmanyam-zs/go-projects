package gofr

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/pkg/gofr/config"
	"developer.zopsmart.com/go/gofr/pkg/log"
)

func TestTraceExporterSuccess(t *testing.T) {
	logger := log.NewMockLogger(io.Discard)
	cfg := config.NewGoDotEnvProvider(logger, "../../configs")
	testcases := struct {
		// exporter input
		exporter string
		host     string
		port     string
		appName  string
	}{cfg.Get("TRACER_EXPORTER"), cfg.Get("TRACER_HOST"), cfg.Get("TRACER_PORT"), "gofr"}

	tp := TraceProvider(testcases.appName, testcases.exporter, testcases.host, testcases.port, logger, cfg)

	assert.NotNil(t, tp, "Failed.\tExpected NotNil Got Nil")
}

func TestTraceExporterFailure(t *testing.T) {
	testcases := []struct {
		// exporter input
		exporter string
		host     string
		port     string
		appName  string
	}{
		{"not zipkin", "localhost", "2005", "gofr"},
		{"zipkin", "localhost", "asd", "gofr"},
		{"gcp", "fakeproject", "0", "sample-api"},
	}

	for _, v := range testcases {
		tp := TraceProvider(v.appName, v.exporter, v.host, v.port, log.NewMockLogger(io.Discard), &config.MockConfig{})

		assert.Nil(t, tp, "Failed.\tExpected Nil Got NotNil")
	}
}
