package gofr

import (
	"os"

	"developer.zopsmart.com/go/gofr/pkg/log"

	"go.opentelemetry.io/otel/trace"
)

type cmdApp struct {
	Router  CMDRouter
	server  *server
	context *Context
}

func (app *cmdApp) Start(logger log.Logger) {
	args := os.Args[1:] // 1st one is the command name itself.
	command := ""

	for _, a := range args {
		if a[1] != '-' {
			command = command + " " + a
		}
	}

	// starts the HTTP server which is used for metrics and healthCheck endpoints.
	go app.server.Start(app.context.Logger)

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

	trace.SpanFromContext(app.context).End()
}
