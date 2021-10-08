package http

import (
	"developer.zopsmart.com/go/gofr/examples/sample-grpc/handler/grpc"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func Get(ctx *gofr.Context) (interface{}, error) {
	if ctx.Param("id") == "1" {
		resp := grpc.Response{
			FirstName:  "Henry",
			SecondName: "Marc",
		}

		return &resp, nil
	}

	return nil, errors.EntityNotFound{Entity: "name", ID: "2"}
}
