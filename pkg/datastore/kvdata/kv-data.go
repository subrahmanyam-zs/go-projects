package kvdata

import (
	"context"
	"encoding/json"
	"net/http"

	"developer.zopsmart.com/go/gofr/pkg/errors"
)

type Config struct {
	URL       string
	AppKey    string
	SharedKey string
}

type client struct {
	kvDataSvc HTTPService
}

//nolint:revive //client should not be accessed directly
func New(kvDataSvc HTTPService) client {
	return client{kvDataSvc: kvDataSvc}
}

func (c client) Get(ctx context.Context, key string) (string, error) {
	resp, err := c.kvDataSvc.Get(ctx, "data/"+key, nil)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", c.getError(resp.Body, resp.StatusCode)
	}

	var res = struct {
		Data map[string]string
	}{}

	err = c.bind(resp.Body, &res)
	if err != nil {
		return "", err
	}

	return res.Data[key], nil
}

func (c client) Set(ctx context.Context, key, value string) error {
	input := make(map[string]string)
	input[key] = value

	body, err := json.Marshal(&input)
	if err != nil {
		return err
	}

	resp, err := c.kvDataSvc.Post(ctx, "data", nil, body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return c.getError(resp.Body, resp.StatusCode)
	}

	return nil
}

func (c client) Delete(ctx context.Context, key string) error {
	resp, err := c.kvDataSvc.Delete(ctx, "data/"+key, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return c.getError(resp.Body, resp.StatusCode)
	}

	return nil
}

// getError unmarshalls the error response and returns it.
// If error occurs while unmarshalling, it returns Bind Error.
func (c client) getError(body []byte, statusCode int) error {
	resp := struct {
		Errors []errors.Response `json:"errors"`
	}{}

	if err := c.bind(body, &resp); err != nil {
		return err
	}

	err := errors.MultipleErrors{StatusCode: statusCode}
	for i := range resp.Errors {
		err.Errors = append(err.Errors, &resp.Errors[i])
	}

	return err
}

// bind unmarshalls response body to data and returns Bind Error if an error occurs.
func (c client) bind(body []byte, data interface{}) error {
	err := c.kvDataSvc.Bind(body, data)
	if err != nil {
		return &errors.Response{
			Code:   "Bind Error",
			Reason: "failed to bind response",
			Detail: err,
		}
	}

	return nil
}
