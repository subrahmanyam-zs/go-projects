package grpc

import (
	"context"

	"developer.zopsmart.com/go/gofr/pkg/errors"
)

type Handler struct{}

func (h Handler) Get(ctx context.Context, id *ID) (*Response, error) {
	if id.Id == "1" {
		resp := &Response{
			FirstName:  "First",
			SecondName: "Second",
		}

		return resp, nil
	}

	return nil, errors.EntityNotFound{Entity: "name", ID: id.Id}
}
