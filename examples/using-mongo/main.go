package main

import (
	"developer.zopsmart.com/go/gofr/examples/using-mongo/handlers"
	"developer.zopsmart.com/go/gofr/examples/using-mongo/stores/customer"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	// create the application object
	app := gofr.New()
	app.EnableSwaggerUI()

	store := customer.New()
	h := handlers.New(store)

	// specifying the different routes supported by this service
	app.GET("/customer", h.Get)
	app.POST("/customer", h.Create)
	app.DELETE("/customer", h.Delete)
	app.Server.HTTP.Port = 9097

	app.Start()
}
