package user

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/examples/using-http-service/models"
	"developer.zopsmart.com/go/gofr/examples/using-http-service/services"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
	"developer.zopsmart.com/go/gofr/pkg/gofr/responder"
)

func Test_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := services.NewMockUser(ctrl)

	mockService.EXPECT().Get(gomock.Any(), "Vikash").Return(models.User{Name: "Vikash", Company: "ZopSmart"}, nil)
	mockService.EXPECT().Get(gomock.Any(), "ABC").Return(models.User{}, errors.EntityNotFound{Entity: "User", ID: "ABC"})

	tests := []struct {
		desc string
		name string
		resp interface{}
		err  error
	}{
		{"get with missing params", "", nil, errors.MissingParam{Param: []string{"name"}}},
		{"get succuss", "Vikash", models.User{Name: "Vikash", Company: "ZopSmart"}, nil},
		{"get non existent entity", "ABC", nil, errors.EntityNotFound{Entity: "User", ID: "ABC"}},
	}

	for i, tc := range tests {
		req := httptest.NewRequest(http.MethodGet, "http://dummy", nil)
		ctx := gofr.NewContext(responder.NewContextualResponder(httptest.NewRecorder(), req), request.NewHTTPRequest(req), nil)

		h := New(mockService)

		ctx.SetPathParams(map[string]string{"name": tc.name})
		resp, err := h.Get(ctx)

		assert.Equal(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)

		assert.Equal(t, tc.resp, resp, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}
