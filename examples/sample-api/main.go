package main

import (
	"developer.zopsmart.com/go/gofr/examples/sample-api/handler"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	// create the application object
	app := gofr.New()

	// add a handler
	app.GET("/hello-world", handler.HelloWorld)

	// handler can access the parameters from context.
	app.GET("/hello", handler.HelloName)

	// handler function can send response in JSON
	app.GET("/json", handler.JSONHandler)

	// handler returns response based on PathParam
	app.GET("/user/{name}", handler.UserHandler)

	// Handler function which throws error
	app.GET("/error", handler.ErrorHandler)

	// Handler function which uses logging
	app.GET("/log", handler.HelloLogHandler)

	// Handler function return response of type raw
	app.GET("/raw", handler.Raw)

	// start the server
	app.Start()
}
