package user

import (
	"net/http"

	"developer.zopsmart.com/go/gofr/examples/using-http-service/models"
	"developer.zopsmart.com/go/gofr/examples/using-http-service/services"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type service struct {
	svc services.HTTPService
}

// New is factory function for service layer
func New(svc services.HTTPService) services.User {
	return service{svc: svc}
}

func (s service) Get(ctx *gofr.Context, name string) (models.User, error) {
	resp, err := s.svc.Get(ctx, "user/"+name, nil)
	if err != nil {
		return models.User{}, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		body := struct {
			Data models.User `json:"data"`
		}{}

		err = s.svc.Bind(resp.Body, &body)
		if err == nil {
			return body.Data, nil
		}
	default:
		err := errors.MultipleErrors{Errors: []error{&errors.Response{}}}

		e := s.svc.Bind(resp.Body, &err)
		if e == nil {
			return models.User{}, err
		}
	}

	return models.User{}, &errors.Response{
		StatusCode: http.StatusInternalServerError,
		Code:       "BIND_ERROR",
		Reason:     "failed to bind response from sample service",
	}
}
