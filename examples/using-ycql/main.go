package main

import (
	"developer.zopsmart.com/go/gofr/examples/using-ycql/handlers"
	"developer.zopsmart.com/go/gofr/examples/using-ycql/store/shop"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	// Create the application object
	k := gofr.New()

	// initialize the handler
	h := handlers.New(shop.Shop{})
	// added get handler
	k.GET("/shop", h.Get)
	// added create handler
	k.POST("/shop", h.Create)
	// added update handler
	k.PUT("/shop/{id}", h.Update)
	// added delete handler
	k.DELETE("/shop/{id}", h.Delete)
	// server  can  start at custom port
	k.Server.HTTP.Port = 9005

	// server start
	k.Start()
}
