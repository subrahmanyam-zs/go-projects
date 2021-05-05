package gofr

import (
	"strings"

	"contrib.go.opencensus.io/exporter/stackdriver"
	zkExporter "contrib.go.opencensus.io/exporter/zipkin"
	zk "github.com/openzipkin/zipkin-go"
	zkHTTP "github.com/openzipkin/zipkin-go/reporter/http"
	"go.opencensus.io/trace"
)

type exporter struct {
	name    string
	host    string
	port    string
	appName string
}

func TraceExporter(appName, exporterName, exporterHost, exporterPort string) trace.Exporter {
	exporterName = strings.ToLower(exporterName)
	e := exporter{
		name:    exporterName,
		host:    exporterHost,
		port:    exporterPort,
		appName: appName,
	}

	switch exporterName {
	case "zipkin":
		return e.getZipkinExporter()
	case "gcp":
		return getGCPExporter(exporterHost)
	default:
		return nil
	}
}

func (e *exporter) getZipkinExporter() trace.Exporter {
	localEndpoint, err := zk.NewEndpoint(e.appName, e.host)
	if err != nil {
		return nil
	}

	url := "http://" + e.host + ":" + e.port + "/api/v2/spans"
	reporter := zkHTTP.NewReporter(url)
	ze := zkExporter.NewExporter(reporter, localEndpoint)

	return ze
}

func getGCPExporter(projectID string) trace.Exporter {
	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: projectID,
	})
	if err != nil {
		return nil
	}

	return exporter
}
