package main

import (
	"github.com/zopsmart/gofr/examples/using-eventhub/handlers"
	"github.com/zopsmart/gofr/pkg/gofr"
)

func main() {
	k := gofr.New()

	k.GET("/pub", handlers.Producer)
	k.GET("/sub", handlers.Consumer)

	k.Server.HTTP.Port = 9113
	k.Start()
}
