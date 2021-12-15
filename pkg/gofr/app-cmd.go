package gofr

import (
	"net/http"
	"os"

	"go.opencensus.io/trace"

	"developer.zopsmart.com/go/gofr/pkg/log"
)

type cmdApp struct {
	Router         CMDRouter
	metricSvr      *metricServer
	healthCheckSvr *healthCheckServer
	context        *Context
	tracingSpan    *trace.Span
}

type metricServer struct {
	server *http.Server
	port   int
	route  string
}

type healthCheckServer struct {
	server *http.Server
	port   int
	route  string
}

func (app *cmdApp) Start(logger log.Logger) {
	args := os.Args[1:] // 1st one is the command name itself.
	command := ""

	for _, a := range args {
		if a[1] != '-' {
			command = command + " " + a
		}
	}

	// start the metric server
	app.metricSvr.server = metricsServer(logger, app.metricSvr.port, app.metricSvr.route)

	// start the health-check server
	go func() {
		app.context.Logger.Infof("Starting health-check server at :%v", app.healthCheckSvr.port)

		err := app.healthCheckSvr.server.ListenAndServe()
		if err != nil {
			app.context.Logger.Errorf("error in health-check server %v", err)
		}
	}()

	h := app.Router.handler(command)
	if h == nil {
		app.context.resp.Respond("No Command Found!", nil)
		return
	}

	data, err := h(app.context)
	if err != nil {
		app.context.resp.Respond(nil, err)
	} else {
		app.context.resp.Respond(data, nil)
	}

	app.tracingSpan.End()
}
