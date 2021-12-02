package gofr

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

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
	app.healthCheckSvr.server = app.healthCheckHandler(logger, app.healthCheckSvr.port, app.healthCheckSvr.route)

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

func (app *cmdApp) healthCheckHandler(logger log.Logger, port int, route string) *http.Server {
	mux := http.NewServeMux()
	healthResp, err := HealthHandler(app.context)

	mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if err != nil {
			logger.Error(err)

			data, _ := json.Marshal(err)

			_, err := w.Write(data)
			if err != nil {
				logger.Error(err)
				return
			}
		} else {
			data, _ := json.Marshal(healthResp)
			_, err := w.Write(data)
			if err != nil {
				logger.Error(err)
				return
			}
		}
	})

	var srv = &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: mux,
	}

	logger.Infof("Starting health-check server at :%v", port)

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			logger.Errorf("error in health-check server %v", err)
		}
	}()

	return srv
}
