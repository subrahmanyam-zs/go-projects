package gofr

import (
	"bytes"
	"reflect"
	"testing"

	"developer.zopsmart.com/go/gofr/pkg/log"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/sdk/trace"
)

func TestTraceExporterSuccess(t *testing.T) {
	testcases := []struct {
		// exporter input
		name    string
		host    string
		port    string
		appName string
	}{
		{"zipkin", "invalid", "2005", "gofr"},
	}

	for _, v := range testcases {
		b := new(bytes.Buffer)
		logger := log.NewMockLogger(b)
		tp := TraceProvider(v.appName, v.name, v.host, v.port, logger)

		if assert.Nil(t, tp) {
			t.Errorf("Failed.\tExpected NotNil Got Nil")
		}
	}
}

func TestTraceExporterFailure(t *testing.T) {
	testcases := []struct {
		// exporter input
		name    string
		host    string
		port    string
		appName string
		tp      *trace.TracerProvider
	}{
		{"not zipkin", "localhost", "2005", "gofr", nil},
		{"gcp", "fakeproject", "0", "gofr", nil},
	}

	for _, v := range testcases {
		b := new(bytes.Buffer)
		logger := log.NewMockLogger(b)
		tp := TraceProvider(v.appName, v.name, v.host, v.port, logger)

		if !reflect.DeepEqual(tp, v.tp) {
			t.Errorf("Failed.\tExpected %v\tGot %v\n", v.tp, tp)
		}
	}
}
