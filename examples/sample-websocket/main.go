package main

import (
	"developer.zopsmart.com/go/gofr/examples/sample-websocket/handlers"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	app := gofr.New()

	app.GET("/", handlers.HomeHandler)
	app.GET("/ws", handlers.WSHandler)

	app.Server.WSUpgrader.WriteBufferSize = 4096

	app.Start()
}
