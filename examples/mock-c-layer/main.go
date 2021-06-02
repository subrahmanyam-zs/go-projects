package main

import (
	"developer.zopsmart.com/go/gofr/examples/mock-c-layer/handler"
	"developer.zopsmart.com/go/gofr/examples/mock-c-layer/store/brand"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	k := gofr.New()

	// initialize the core and consumer layers
	brandCore := brand.New()
	brandConsumer := handler.New(brandCore)

	// Specifying the different routes supported by this service
	k.GET("/brand", brandConsumer.Get)
	k.POST("/brand", brandConsumer.Create)

	k.Start()
}
