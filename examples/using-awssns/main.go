package main

import (
	"developer.zopsmart.com/go/gofr/examples/using-awssns/handlers"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	app := gofr.New()

	app.Server.ValidateHeaders = false

	app.POST("/publish", handlers.Publisher)
	app.GET("/subscribe", handlers.Subscriber)

	app.Start()
}
