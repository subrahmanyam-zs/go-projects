package handlers

import (
	"time"

	"developer.zopsmart.com/go/gofr/pkg/datastore/pubsub"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/types"
)

type Person struct {
	ID    string `avro:"Id"`
	Name  string `avro:"Name"`
	Email string `avro:"Email"`
}

func Producer(c *gofr.Context) (interface{}, error) {
	id := c.Param("id")
	start := time.Now()

	err := c.PublishEvent("", Person{
		ID:    id,
		Name:  "Rohan",
		Email: "rohan@email.xyz",
	}, map[string]string{"test": "test"})
	if err != nil {
		return nil, err
	}

	err = c.Metric.ObserveHistogram(PublishEventHistogram, time.Since(start).Seconds())

	return nil, err
}

func Consumer(c *gofr.Context) (interface{}, error) {
	p := Person{}
	start := time.Now()

	message, err := c.Subscribe(&p)
	if err != nil {
		return nil, err
	}

	err = c.Metric.ObserveSummary(ConsumeEventSummary, time.Since(start).Seconds())

	return types.Response{Data: p, Meta: message}, err
}

func ConsumerWithCommit(c *gofr.Context) (interface{}, error) {
	p := Person{}

	count := 0
	message, err := c.SubscribeWithCommit(func(message *pubsub.Message) (bool, bool) {
		count++
		c.Logger.Infof("Consumed %v message(s), offset: %v, topic: %v", count, message.Offset, message.Topic)

		for count <= 2 {
			return true, true
		}

		for count <= 5 {
			return false, true
		}

		return false, false
	})

	if err != nil {
		return nil, err
	}

	return types.Response{Data: p, Meta: message}, nil
}
