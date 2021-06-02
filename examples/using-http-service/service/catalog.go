package service

import (
	"context"

	"developer.zopsmart.com/go/gofr/pkg/log"
	"developer.zopsmart.com/go/gofr/pkg/service"
)

type CatalogService interface {
	GetBrandByID(ctx context.Context, id string) interface{}
}

type catalog struct {
	httpService service.HTTP
}

func New(url string, logger log.Logger) catalog {
	httpSvc := service.NewHTTPServiceWithOptions(url, logger, nil)
	return catalog{httpService: httpSvc}
}

func (c catalog) GetBrandByID(ctx context.Context, id string) interface{} {
	c.httpService.PropagateHeaders()
	svcResp, err := c.httpService.Get(ctx, id, nil)
	if err != nil {
		return nil
	}

	resp := make(map[string]interface{})

	err = c.httpService.Bind(svcResp.Body, &resp)
	if err != nil {
		return nil
	}

	return resp
}
