package gofr

import (
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/pkg/gofr/config"
	"developer.zopsmart.com/go/gofr/pkg/log"
)

func TestTraceExporterSuccess(t *testing.T) {
	cfg := config.NewGoDotEnvProvider(log.NewMockLogger(io.Discard), "../../configs")
	err := tracerProvider(cfg)

	assert.Nil(t, err, "Failed.\tExpected NotNil Got Nil")
}

func TestTraceExporterFailure(t *testing.T) {
	testcases := []struct {
		// exporter input
		exporter string
		host     string
		port     string
		appName  string
	}{
		{"not zipkin", "localhost", "9411", "gofr"},
		{"zipkin", "localhost", "asd", "gofr"},
		{"gcp", "fakeproject", "0", "sample-api"},
	}

	for _, v := range testcases {
		tracerUrl := fmt.Sprintf("http://%v:%v", v.host, v.port)
		cfg := &config.MockConfig{Data: map[string]string{
			"TRACER_EXPORTER": v.exporter,
			"TRACER_HOST":     v.host,
			"TRACER_PORT":     v.port,
			"TRACER_URL":      tracerUrl,
		}}

		err := tracerProvider(cfg)

		assert.NotNil(t, err, "Failed.\tExpected Nil Got NotNil")
	}
}
