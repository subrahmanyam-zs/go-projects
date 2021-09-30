package main

import (
	"developer.zopsmart.com/go/gofr/examples/using-http2/handler"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	app := gofr.New()

	// add a handler
	app.GET("/static/{name}", handler.ServeStatic)
	app.GET("/home", handler.HomeHandler)

	// set https port and redirect
	app.Server.HTTPS.Port = 1449
	app.Server.HTTP.RedirectToHTTPS = false

	// http port
	app.Server.HTTP.Port = 9017

	app.Start()
}
