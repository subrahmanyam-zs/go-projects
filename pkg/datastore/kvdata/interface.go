package kvdata

import (
	"context"

	"developer.zopsmart.com/go/gofr/pkg/service"
)

type HTTPService interface {
	Get(ctx context.Context, api string, params map[string]interface{}) (*service.Response, error)
	Post(ctx context.Context, api string, params map[string]interface{}, body []byte) (*service.Response, error)
	Delete(ctx context.Context, api string, body []byte) (*service.Response, error)
	Bind(resp []byte, i interface{}) error
}
