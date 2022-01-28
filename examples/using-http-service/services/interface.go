package services

import (
	"context"

	"developer.zopsmart.com/go/gofr/examples/using-http-service/models"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/service"
)

type HTTPService interface {
	Get(ctx context.Context, api string, params map[string]interface{}) (*service.Response, error)
	Bind(resp []byte, i interface{}) error
}

type User interface {
	Get(ctx *gofr.Context, name string) (models.User, error)
}
