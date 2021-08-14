package gofr

import (
	"io"
	"testing"

	"developer.zopsmart.com/go/gofr/pkg/gofr/config"
	"developer.zopsmart.com/go/gofr/pkg/log"

	"github.com/stretchr/testify/assert"
)

func TestTraceExporterSuccess(t *testing.T) {
	testcases := struct {
		// exporter input
		name    string
		host    string
		port    string
		appName string
	}{"zipkin", "localhost", "2005", "gofr"}

	tp := TraceProvider(testcases.appName, testcases.name, testcases.host, testcases.port, log.NewMockLogger(io.Discard), &config.MockConfig{})

	assert.NotNil(t, tp, "Failed.\tExpected NotNil Got Nil")
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
		{"zipkin", "localhost", "asd", "gofr"},
		{"gcp", "fakeproject", "0", "sample-api"},
	}

	for _, v := range testcases {
		tp := TraceProvider(v.appName, v.name, v.host, v.port, log.NewMockLogger(io.Discard), &config.MockConfig{})

		assert.Nil(t, tp, "Failed.\tExpected Nil Got NotNil")
	}
}
