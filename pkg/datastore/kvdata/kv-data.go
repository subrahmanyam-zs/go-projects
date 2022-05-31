package kvdata

import (
	"context"
	"encoding/json"
	"net/http"

	"developer.zopsmart.com/go/gofr/pkg/errors"
)

//nolint
type KvDataConfig struct {
	URL       string
	AppKey    string
	SharedKey string
}

type kvDataClient struct {
	kvDataSvc HTTPService
}

//nolint
func New(kvDataSvc HTTPService) kvDataClient {
	return kvDataClient{kvDataSvc: kvDataSvc}
}

func (k kvDataClient) Get(ctx context.Context, key string) (string, error) {
	resp, err := k.kvDataSvc.Get(ctx, "data/"+key, nil)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", k.getError(resp.Body, resp.StatusCode)
	}

	var res = struct {
		Data map[string]string
	}{}

	err = k.bind(resp.Body, &res)
	if err != nil {
		return "", err
	}

	return res.Data[key], nil
}

func (k kvDataClient) Set(ctx context.Context, key, value string) error {
	input := make(map[string]string)
	input[key] = value

	body, err := json.Marshal(&input)
	if err != nil {
		return err
	}

	resp, err := k.kvDataSvc.Post(ctx, "data", nil, body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return k.getError(resp.Body, resp.StatusCode)
	}

	return nil
}

func (k kvDataClient) Delete(ctx context.Context, key string) error {
	resp, err := k.kvDataSvc.Delete(ctx, "data/"+key, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return k.getError(resp.Body, resp.StatusCode)
	}

	return nil
}

// getError unmarshalls the error response and returns it.
// If error occurs while unmarshalling, it returns Bind Error.
func (k kvDataClient) getError(body []byte, statusCode int) error {
	resp := struct {
		Errors []errors.Response `json:"errors"`
	}{}

	if err := k.bind(body, &resp); err != nil {
		return err
	}

	err := errors.MultipleErrors{StatusCode: statusCode}
	for i := range resp.Errors {
		err.Errors = append(err.Errors, &resp.Errors[i])
	}

	return err
}

// bind unmarshalls response body to data and returns Bind Error if an error occurs.
func (k kvDataClient) bind(body []byte, data interface{}) error {
	err := k.kvDataSvc.Bind(body, data)
	if err != nil {
		return &errors.Response{
			Code:   "Bind Error",
			Reason: "failed to bind response",
			Detail: err,
		}
	}

	return nil
}
