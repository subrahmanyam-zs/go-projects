package main

import (
	"developer.zopsmart.com/go/gofr/examples/using-redis/handler"
	"developer.zopsmart.com/go/gofr/examples/using-redis/store"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	// Create the application object
	k := gofr.New()

	s := store.New()
	h := handler.New(s)

	// Specifying the different routes supported by this service
	k.GET("/config/{key}", h.GetKey)
	k.POST("/config", h.SetKey)
	k.DELETE("/config/{key}", h.DeleteKey)

	// Starting the server on a custom port
	k.Server.HTTP.Port = 9091
	k.Start()
}
