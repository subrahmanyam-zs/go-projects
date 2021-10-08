package main

import (
	"developer.zopsmart.com/go/gofr/examples/using-eventhub/handlers"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	app := gofr.New()

	app.GET("/pub", handlers.Producer)
	app.GET("/sub", handlers.Consumer)

	app.Server.HTTP.Port = 9113
	app.Start()
}
