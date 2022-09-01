package main

import (
	handlers "developer.zopsmart.com/go/gofr/examples/data-layer-with-cassandra/handlers/person"
	stores "developer.zopsmart.com/go/gofr/examples/data-layer-with-cassandra/stores/person"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	// create the application object
	app := gofr.New()

	app.Server.ValidateHeaders = false

	s := stores.New()
	h := handlers.New(s)
	// add get handler
	app.GET("/persons", h.Get)
	// add post handler
	app.POST("/persons", h.Create)
	// add a delete handler
	app.DELETE("/persons/{id}", h.Delete)
	// add a put handler
	app.PUT("/persons/{id}", h.Update)

	// starting the server on a custom port
	app.Server.HTTP.Port = 9094
	app.Server.MetricsPort = 2123
	// start the server
	app.Start()
}
