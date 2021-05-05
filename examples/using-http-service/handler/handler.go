package handler

import (
	"github.com/zopsmart/gofr/examples/using-http-service/service"
	"github.com/zopsmart/gofr/pkg/errors"
	"github.com/zopsmart/gofr/pkg/gofr"
)

type Handler struct {
	catalogService service.CatalogService
}

func New(catalogSvc service.CatalogService) Handler {
	return Handler{catalogService: catalogSvc}
}

// Get is a handler function of type gofr.Handler that uses HTTP Service to make downstream calls
func (h Handler) Get(c *gofr.Context) (interface{}, error) {
	id := c.PathParam("id")
	if id == "" {
		return nil, errors.MissingParam{Param: []string{"id"}}
	}
	return h.catalogService.GetBrandByID(c, id), nil
}
