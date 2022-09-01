package main

import (
	"developer.zopsmart.com/go/gofr/examples/data-layer-with-postgres/handler"
	"developer.zopsmart.com/go/gofr/examples/data-layer-with-postgres/store"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	app := gofr.New()

	s := store.New()
	h := handler.New(s)

	// specifying the different routes supported by this service
	app.GET("/customer", h.Get)
	app.GET("/customer/{id}", h.GetByID)
	app.POST("/customer", h.Create)
	app.PUT("/customer/{id}", h.Update)
	app.DELETE("/customer/{id}", h.Delete)

	// starting the server on a custom port
	app.Server.HTTP.Port = 9092
	app.Server.MetricsPort = 2325
	app.Start()
}
