package main

import (
	"github.com/zopsmart/gofr/examples/sample-websocket/handlers"
	"github.com/zopsmart/gofr/pkg/gofr"
)

func main() {
	k := gofr.New()

	k.GET("/", handlers.HomeHandler)
	k.GET("/ws", handlers.WSHandler)

	k.Server.WSUpgrader.WriteBufferSize = 4096

	k.Start()
}
