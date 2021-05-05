package main

import (
	"github.com/zopsmart/gofr/examples/mock-c-layer/handler"
	"github.com/zopsmart/gofr/examples/mock-c-layer/store/brand"
	"github.com/zopsmart/gofr/pkg/gofr"
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
