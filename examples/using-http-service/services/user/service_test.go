package user

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/examples/using-http-service/models"
	"developer.zopsmart.com/go/gofr/examples/using-http-service/services"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/log"
	svc "developer.zopsmart.com/go/gofr/pkg/service"
)

func initializeTest(t *testing.T) (services.User, mockHTTPService, *gofr.Context) {
	ctrl := gomock.NewController(t)
	mock := services.NewMockHTTPService(ctrl)
	mockSvc := mockHTTPService{MockHTTPService: mock}
	s := New(mockSvc)

	g := gofr.Gofr{Logger: log.NewMockLogger(io.Discard)}
	ctx := gofr.NewContext(nil, nil, &g)

	return s, mockSvc, ctx
}

func TestService_Get(t *testing.T) {
	s, mockSvc, ctx := initializeTest(t)

	var (
		serviceErr = errors.Error("error from service call")

		serverErr = errors.MultipleErrors{StatusCode: http.StatusInternalServerError, Errors: []error{&errors.Response{
			Code:     "Internal Server Error",
			Reason:   "connection timed out",
			DateTime: errors.DateTime{Value: "2021-11-03T11:01:13.124Z", TimeZone: "IST"},
		}}}

		errResp = &svc.Response{
			StatusCode: http.StatusInternalServerError,
			Body: []byte(`
		{
			"errors": [
				{
					"code": "Internal Server Error",
					"reason": "connection timed out",
					"datetime": {
						"value": "2021-11-03T11:01:13.124Z",
						"timezone": "IST"
					}
				}
			]
		}`),
		}

		resp = &svc.Response{
			StatusCode: http.StatusOK,
			Body: []byte(`
		{
    		"data": {
        		"name": "Vikash",
        		"company": "ZopSmart"
    		}
		}`),
		}
	)

	tests := []struct {
		desc     string
		mockResp *svc.Response // mock response from get call
		mockErr  error         // mock error from get call
		output   models.User
		err      error
	}{
		{"error from get call", nil, serviceErr, models.User{}, serviceErr},
		{"error response from get call", errResp, nil, models.User{}, serverErr},
		{"success case", resp, nil, models.User{Name: "Vikash", Company: "ZopSmart"}, nil},
	}

	for i, tc := range tests {
		mockSvc.EXPECT().Get(gomock.Any(), "user/Vikash", nil).Return(tc.mockResp, tc.mockErr)

		output, err := s.Get(ctx, "Vikash")

		assert.Equal(t, tc.err, err, "TEST[%d], failed.\n%s", i+1, tc.desc)

		assert.Equal(t, tc.output, output, "TEST[%d], failed.\n%s", i+1, tc.desc)
	}
}

func TestService_GetBindError(t *testing.T) {
	s, mockSvc, ctx := initializeTest(t)

	var bindErr = &errors.Response{
		Code:   "Bind Error",
		Reason: "failed to bind response",
		Detail: errors.Error("bind error"),
	}

	tests := []struct {
		desc     string
		mockResp *svc.Response // mock response from get call
	}{
		{"error in binding error response", &svc.Response{StatusCode: http.StatusBadRequest, Body: []byte(`invalid body`)}},
		{"error in binding data response", &svc.Response{StatusCode: http.StatusOK, Body: []byte(`invalid body`)}},
	}

	for i, tc := range tests {
		mockSvc.EXPECT().Get(gomock.Any(), "user/Vikash", nil).Return(tc.mockResp, nil)
		mockSvc.EXPECT().Bind(tc.mockResp.Body, gomock.Any()).Return(errors.Error("bind error"))

		output, err := s.Get(ctx, "Vikash")

		assert.Equal(t, bindErr, err, "TEST[%d], failed.\n%s", i, tc.desc)

		assert.Equal(t, models.User{}, output, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

type mockHTTPService struct {
	*services.MockHTTPService // embed mock of HTTPService interface
}

// override the Bind method
func (m mockHTTPService) Bind(resp []byte, i interface{}) error {
	if err := json.Unmarshal(resp, i); err != nil {
		return m.MockHTTPService.Bind(resp, i)
	}

	return nil
}
