package gofr

import (
	"testing"

	"go.opencensus.io/trace"
)

func TestTraceExporter(t *testing.T) {
	testcases := []struct {
		// exporter input
		name    string
		host    string
		port    string
		appName string

		// output
		exporter trace.Exporter
	}{
		{"zipkin", "invalid", "2005", "gofr", nil},
		{"not zipkin", "localhost", "2005", "gofr", nil},
		{"gcp", "fakeproject", "0", "gofr", nil},
	}

	for _, v := range testcases {
		exporter := TraceExporter(v.appName, v.name, v.host, v.port)

		if exporter != v.exporter {
			t.Errorf("Failed.\tExpected %v\tGot %v\n", v.exporter, exporter)
		}
	}
}

func TestGCPTrace(t *testing.T) {
	tests := []struct {
		name      string
		projectID string
		want      trace.Exporter
	}{
		{"exporter creation failed", "", nil},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := getGCPExporter(tt.projectID)
			if got != tt.want {
				t.Errorf("getGCPExporter() = %v, want %v", got, tt.want)
			}
		})
	}
}
