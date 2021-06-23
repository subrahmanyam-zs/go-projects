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

	return nil, c.Notifier.Publish(message)
}

func Subscriber(c *gofr.Context) (interface{}, error) {
	data := map[string]interface{}{}
	msg, err := c.Notifier.SubscribeWithResponse(&data)

	if err != nil {
		return nil, err
	}

	return types.Response{Data: data, Meta: msg}, nil
}
