package user

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/examples/using-http-service/models"
	"developer.zopsmart.com/go/gofr/examples/using-http-service/services"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	httpService "developer.zopsmart.com/go/gofr/pkg/service"
)

type mockBind struct{}

func (m mockBind) Bind(resp []byte, i interface{}) error {
	err := json.NewDecoder(bytes.NewBuffer(resp)).Decode(&i)
	return err
}

func (m mockBind) Get(ctx context.Context, api string, params map[string]interface{}) (*httpService.Response, error) {
	input := []byte(`{"errors":[{"code":"400"}]}`)
	return &httpService.Response{Body: input, StatusCode: http.StatusBadRequest}, nil
}

func Test_Get(t *testing.T) {
	tests := []struct {
		desc string
		resp models.User
		err  error
	}{
		{"call to service.Get throws error", models.User{},
			errors.MultipleErrors{StatusCode: http.StatusInternalServerError, Errors: []error{errors.Error("core error")}}},
		{"call to Bind method throws error", models.User{},
			&errors.Response{StatusCode: http.StatusInternalServerError, Code: "BIND_ERROR", Reason: "failed to bind response from sample service"}},
		{"success case", models.User{Name: "Vikash", Company: "ZopSmart"}, nil},
	}

	for i, tc := range tests {
		h := New(services.New(i))

		ctx := gofr.NewContext(nil, nil, gofr.New())
		resp, err := h.Get(ctx, "Vikash")

		assert.Equal(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)

		assert.Equal(t, tc.resp, resp, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

func Test_DefaultCase(t *testing.T) {
	h := New(mockBind{})

	ctx := gofr.NewContext(nil, nil, gofr.New())

	_, err := h.Get(ctx, "vikash")
	if err == nil {
		t.Errorf("failed expected error got nil")
	}
}
