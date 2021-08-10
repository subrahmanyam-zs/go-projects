package gofr

import (
	"io"
	"testing"

	"developer.zopsmart.com/go/gofr/pkg/log"

	"github.com/stretchr/testify/assert"
)

func TestTraceExporterSuccess(t *testing.T) {
	testcases := []struct {
		// exporter input
		name    string
		host    string
		port    string
		appName string
	}{
		{"zipkin", "localhost", "2005", "gofr"},
	}

	for _, v := range testcases {
		logger := log.NewMockLogger(io.Discard)
		tp := TraceProvider(v.appName, v.name, v.host, v.port, logger)

		assert.NotNil(t, tp, "Failed.\tExpected NotNil Got Nil")
	}
}

func TestTraceExporterFailure(t *testing.T) {
	testcases := []struct {
		// exporter input
		name    string
		host    string
		port    string
		appName string
	}{
		{"not zipkin", "localhost", "2005", "gofr"},
		{"gcp", "fakeproject", "0", "gofr-dev"},
	}

	for _, v := range testcases {
		logger := log.NewMockLogger(io.Discard)
		tp := TraceProvider(v.appName, v.name, v.host, v.port, logger)

		assert.Nil(t, tp, "Failed.\tExpected Nil Got NotNil")
	}
}
