package http

import (
	"github.com/zopsmart/gofr/examples/sample-grpc/handler/grpc"
	"github.com/zopsmart/gofr/pkg/errors"
	"github.com/zopsmart/gofr/pkg/gofr"
)

func Get(c *gofr.Context) (interface{}, error) {
	if c.Param("id") == "1" {
		resp := grpc.Response{
			FirstName:  "Henry",
			SecondName: "Marc",
		}

		return &resp, nil
	}

	return nil, errors.EntityNotFound{Entity: "name", ID: "2"}
}
