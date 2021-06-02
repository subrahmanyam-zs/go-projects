package main

import (
	"developer.zopsmart.com/go/gofr/examples/using-http2/handler"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	k := gofr.New()

	// add a handler
	k.GET("/static/{name}", handler.ServeStatic)
	k.GET("/home", handler.HomeHandler)

	// set https port and redirect
	k.Server.HTTPS.Port = 1449
	k.Server.HTTP.RedirectToHTTPS = false

	// http port
	k.Server.HTTP.Port = 9017

	k.Start()
}
