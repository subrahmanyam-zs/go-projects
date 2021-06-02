package handlers

import (
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

	return nil, c.PublishEvent("", Person{
		ID:    id,
		Name:  "Rohan",
		Email: "rohan@email.xyz",
	}, nil)
}

func Consumer(c *gofr.Context) (interface{}, error) {
	p := Person{}
	message, err := c.Subscribe(&p)

	return types.Response{Data: p, Meta: message}, err
}
