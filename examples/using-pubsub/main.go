package main

import (
	"developer.zopsmart.com/go/gofr/examples/using-pubsub/handlers"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	k := gofr.New()

	k.GET("/pub", handlers.Producer)
	k.GET("/sub", handlers.Consumer)
	k.GET("/subCommit", handlers.ConsumerWithCommit)

	k.Server.HTTP.Port = 9112
	k.Start()
}
