package main

import (
	"developer.zopsmart.com/go/gofr/examples/using-pubsub/handlers"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	k := gofr.New()

	err := k.NewHistogram(handlers.PublishEventHistogram,
		"Histogram for time taken to publish event in seconds",
		[]float64{.001, .003, .005, .01, .025, .05, .1, .2, .3, .4, .5, .75, 1, 2, 3, 5, 10, 30})
	if err != nil {
		k.Logger.Warnf("error while creating histogram, %v", err)
	}

	err = k.NewSummary(handlers.ConsumeEventSummary,
		"Summary for time taken to consume event in seconds")
	if err != nil {
		k.Logger.Warnf("error while creating summary, %v", err)
	}

	k.GET("/pub", handlers.Producer)
	k.GET("/sub", handlers.Consumer)
	k.GET("/subCommit", handlers.ConsumerWithCommit)

	k.Server.HTTP.Port = 9112
	k.Start()
}
