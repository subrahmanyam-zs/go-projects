package main

import (
	"developer.zopsmart.com/go/gofr/examples/using-redis/handler"
	"developer.zopsmart.com/go/gofr/examples/using-redis/store"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	// Create the application object
	k := gofr.New()

	s := store.New()
	h := handler.New(s)

	err := k.NewGauge(handler.ReqContentLengthGauge, "Gauge of content-length of request")
	if err != nil {
		k.Logger.Warnf("error while creating Gauge, %v", err)
	}

	err = k.NewCounter(handler.InvalidBodyCounter, "it does count for invalid request body")
	if err != nil {
		k.Logger.Warnf("error while creating counter, %v", err)
	}

	err = k.NewCounter(handler.NumberOfSetsCounter, "it does count for set requests", "status")
	if err != nil {
		k.Logger.Warnf("error while creating counter, %v", err)
	}

	// Specifying the different routes supported by this service
	k.GET("/config/{key}", h.GetKey)
	k.POST("/config", h.SetKey)
	k.DELETE("/config/{key}", h.DeleteKey)

	k.Start()
}
