package main

import (
	"developer.zopsmart.com/go/gofr/examples/using-awssns/handlers"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	k := gofr.New()

	k.Server.ValidateHeaders = false

	k.POST("/publish", handlers.Publisher)
	k.GET("/subscribe", handlers.Subscriber)

	k.Start()
}
