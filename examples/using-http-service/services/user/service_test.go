package user

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/examples/using-http-service/models"
	"developer.zopsmart.com/go/gofr/examples/using-http-service/services"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func Test_Get(t *testing.T) {
	testcases := []struct {
		description string
		response    models.User
		err         error
	}{
		{"call to service.Get throws error", models.User{}, errors.MultipleErrors{
			StatusCode: http.StatusInternalServerError, Errors: []error{errors.Error("core error")}}},
		{"call to Bind method throws error", models.User{}, &errors.Response{StatusCode: http.StatusInternalServerError,
			Code: "BIND_ERROR", Reason: "failed to bind response from sample service"}},
		{"success case", models.User{Name: "Vikash", Company: "ZopSmart"}, nil},
	}

	for i := range testcases {
		h := New(services.New(i))

		ctx := gofr.NewContext(nil, nil, gofr.New())
		resp, err := h.Get(ctx, "Vikash")

		assert.Equal(t, testcases[i].err, err)

		assert.Equal(t, testcases[i].response, resp)
	}
}
