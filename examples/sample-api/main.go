package main

import (
	"developer.zopsmart.com/go/gofr/examples/sample-api/datastore"
	"developer.zopsmart.com/go/gofr/examples/sample-api/handler"
	handleremp "developer.zopsmart.com/go/gofr/examples/sample-api/handler"
	"developer.zopsmart.com/go/gofr/examples/sample-api/service"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	// create the application object
	app := gofr.New()

	app.Server.ValidateHeaders = false
	// enabling /swagger endpoint for Swagger UI
	app.EnableSwaggerUI()

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

	app.GET("/hard", handler.Hard)

	empStore := datastore.New()
	empService := service.New(empStore)
	empHandle := handleremp.New(empService)

	//
	app.POST("/employee", empHandle.Post)

	//
	app.PUT("/employee/{id}", empHandle.Put)

	//
	app.DELETE("/employee/{id}", empHandle.Delete)

	//
	app.GET("/employee/{id}", empHandle.Get)

	//
	app.GET("/employee", empHandle.GetAll)
	// start the server
	app.Start()
}
