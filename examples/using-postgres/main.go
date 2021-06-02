package main

import (
	"developer.zopsmart.com/go/gofr/examples/using-postgres/handler"
	"developer.zopsmart.com/go/gofr/examples/using-postgres/store"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	k := gofr.New()

	s := store.New()
	h := handler.New(s)

	// specifying the different routes supported by this service
	k.GET("/customer", h.Get)
	k.GET("/customer/{id}", h.GetByID)
	k.POST("/customer", h.Create)
	k.PUT("/customer/{id}", h.Update)
	k.DELETE("/customer/{id}", h.Delete)

	// starting the server on a custom port
	k.Server.HTTP.Port = 9092
	k.Server.MetricsPort = 2325
	k.Start()
}
