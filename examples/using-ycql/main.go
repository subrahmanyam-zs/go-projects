package main

import (
	handler "developer.zopsmart.com/go/gofr/examples/using-ycql/handlers/shop"
	store "developer.zopsmart.com/go/gofr/examples/using-ycql/stores/shop"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	// Create the application object
	app := gofr.New()
	// initialize store dependency
	s := store.New()
	// initialize the handler
	h := handler.New(s)
	// added get handler
	app.GET("/shop", h.Get)
	// added create handler
	app.POST("/shop", h.Create)
	// added update handler
	app.PUT("/shop/{id}", h.Update)
	// added delete handler
	app.DELETE("/shop/{id}", h.Delete)
	// server  can  start at custom port
	app.Server.HTTP.Port = 9005

	// server start
	app.Start()
}
