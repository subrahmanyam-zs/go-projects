package main

import (
	"developer.zopsmart.com/go/gofr/examples/using-redis/handler"
	"developer.zopsmart.com/go/gofr/examples/using-redis/store"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	// Create the application object
	app := gofr.New()

	s := store.New()
	h := handler.New(s)

	err := app.NewGauge(handler.ReqContentLengthGauge, "Gauge of content-length of request")
	if err != nil {
		app.Logger.Warnf("error while creating Gauge, %v", err)
	}

	err = app.NewCounter(handler.InvalidBodyCounter, "it does count for invalid request body")
	if err != nil {
		app.Logger.Warnf("error while creating counter, %v", err)
	}

	err = app.NewCounter(handler.NumberOfSetsCounter, "it does count for set requests", "status")
	if err != nil {
		app.Logger.Warnf("error while creating counter, %v", err)
	}

	// Specifying the different routes supported by this service
	app.GET("/config/{key}", h.GetKey)
	app.POST("/config", h.SetKey)
	app.DELETE("/config/{key}", h.DeleteKey)

	app.Start()
}
