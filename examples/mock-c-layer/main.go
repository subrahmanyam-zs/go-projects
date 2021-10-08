package main

import (
	"developer.zopsmart.com/go/gofr/examples/mock-c-layer/handler"
	"developer.zopsmart.com/go/gofr/examples/mock-c-layer/store/brand"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	app := gofr.New()

	// initialize the store and handler layers
	store := brand.New()
	h := handler.New(store)

	// Specifying the different routes supported by this service
	app.GET("/brand", h.Get)
	app.POST("/brand", h.Create)

	app.Start()
}
