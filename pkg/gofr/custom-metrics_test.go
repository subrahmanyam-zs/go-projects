package gofr

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/pkg/errors"
)

var metricErr = errors.Error("invalid/duplicate metrics collector registration attempted")

func Test_New(t *testing.T) {

	testcases := []struct {
		desc string
		err  error
	}{
		{"success-case", nil},
		{"error-case", metricErr},
	}

	k := New()

	for i, tc := range testcases {
		err := k.NewCounter("new_counter", "New Counter", "id")
		assert.Equal(t, tc.err, err, "TESTCASE[%v] NewCounter", i)

		err = k.NewHistogram("new_histogram", "New Histogram", []float64{.5, 1, 2}, "id")
		assert.Equal(t, tc.err, err, "TESTCASE[%v] NewHistogram", i)

		err = k.NewGauge("new_gauge", "New Gauge")
		assert.Equal(t, tc.err, err, "TESTCASE[%v] NewGauge", i)

		err = k.NewSummary("new_summary", "New Summary")
		assert.Equal(t, tc.err, err, "TESTCASE[%v] NewSummary", i)
	}
}
