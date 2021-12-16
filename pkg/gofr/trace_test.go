package gofr

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/pkg/gofr/config"
	"developer.zopsmart.com/go/gofr/pkg/log"
)

func TestTraceExporterSuccess(t *testing.T) {
	cfg := config.NewGoDotEnvProvider(log.NewMockLogger(io.Discard), "../../configs")
	err := tracerProvider(cfg)

	assert.NoError(t, err)
}

func TestTraceExporterFailure(t *testing.T) {
	testcases := []struct {
		// exporter input
		exporter string
		url      string
		appName  string
	}{
		{"not zipkin", "http://localhost/9411", "gofr"},
		{"zipkin", "invalid url", "gofr"},
		{"gcp", "http://fakeProject/9411", "sample-api"},
	}

	for i, tc := range testcases {
		cfg := &config.MockConfig{Data: map[string]string{
			"TRACER_EXPORTER": tc.exporter,
			"TRACER_URL":      tc.url,
		}}

		err := tracerProvider(cfg)

		assert.Error(t, err, "Failed[%v]", i)
	}
}
