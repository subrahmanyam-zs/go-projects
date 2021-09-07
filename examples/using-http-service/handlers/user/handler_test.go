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

	testcases := []struct {
		name     string
		response interface{}
		err      error
	}{
		{"", nil, errors.MissingParam{Param: []string{"name"}}},
		{"Vikash", models.User{Name: "Vikash", Company: "ZopSmart"}, nil},
		{"ABC", nil, errors.EntityNotFound{Entity: "User", ID: "ABC"}},
	}

	for i := range testcases {
		req := httptest.NewRequest(http.MethodGet, "http://dummy", nil)
		c := gofr.NewContext(responder.NewContextualResponder(httptest.NewRecorder(), req), request.NewHTTPRequest(req), nil)

		h := New(mockService)

		c.SetPathParams(map[string]string{"name": testcases[i].name})
		resp, err := h.Get(c)

		assert.Equal(t, testcases[i].err, err, "TEST[%d], failed.", i)

		assert.Equal(t, testcases[i].response, resp, "TEST[%d], failed.", i)
	}
}
