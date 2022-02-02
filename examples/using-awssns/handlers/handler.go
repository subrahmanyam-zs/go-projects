package handlers

import (
	"developer.zopsmart.com/go/gofr/examples/using-awssns/entity"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/types"
)

func Publisher(c *gofr.Context) (interface{}, error) {
	var message *entity.Message

	err := c.Bind(&message)
	if err != nil {
		return nil, errors.EntityNotFound{}
	}

	attr := map[string]interface{}{
		"email":   "test@abc.com",
		"version": 1.1,
		"key":     []interface{}{1, 1.999, "value"},
	}

	return nil, c.Notifier.Publish(message, attr)
}

func Subscriber(c *gofr.Context) (interface{}, error) {
	data := map[string]interface{}{}
	msg, err := c.Notifier.SubscribeWithResponse(&data)

	if err != nil {
		return nil, err
	}

	return types.Response{Data: data, Meta: msg}, nil
}
