package main

import (
	"github.com/zopsmart/gofr/examples/using-cassandra/handlers"
	"github.com/zopsmart/gofr/examples/using-cassandra/store/person"
	"github.com/zopsmart/gofr/pkg/gofr"
)

func main() {
	// create the application object
	k := gofr.New()
	k.Server.ValidateHeaders = false

	h := handlers.New(person.Person{})
	// add get handler
	k.GET("/persons", h.Get)
	// add post handler
	k.POST("/persons", h.Create)
	// add a delete handler
	k.DELETE("/persons/{id}", h.Delete)
	// add a put handler
	k.PUT("/persons/{id}", h.Update)

	// starting the server on a custom port
	k.Server.HTTP.Port = 9094
	k.Server.MetricsPort = 2123
	// start the server
	k.Start()
}
