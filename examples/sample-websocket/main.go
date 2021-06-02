package main

import (
	"developer.zopsmart.com/go/gofr/examples/sample-websocket/handlers"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	k := gofr.New()

	k.GET("/", handlers.HomeHandler)
	k.GET("/ws", handlers.WSHandler)

	k.Server.WSUpgrader.WriteBufferSize = 4096

	k.Start()
}
