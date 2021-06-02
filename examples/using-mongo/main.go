package main

import (
	"developer.zopsmart.com/go/gofr/examples/using-mongo/handlers"
	"developer.zopsmart.com/go/gofr/examples/using-mongo/store/customer"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	// create the application object
	k := gofr.New()
	h := handlers.New(customer.Customer{})

	// specifying the different routes supported by this service
	k.GET("/customer", h.Get)
	k.POST("/customer", h.Create)
	k.DELETE("/customer", h.Delete)
	k.Server.HTTP.Port = 9097

	k.Start()
}
