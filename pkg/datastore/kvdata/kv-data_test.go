package kvdata

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/datastore/kvdata/mocks"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/service"
)

const (
	serviceErr = errors.Error("error from service call")
)

func initializeTest(t *testing.T) (datastore.KVData, mockHTTPService, context.Context) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockHTTPService(ctrl)
	mockSvc := mockHTTPService{MockHTTPService: mock}

	s := New(mockSvc)

	ctx := context.Background()

	return s, mockSvc, ctx
}

func TestKVData_Set(t *testing.T) {
	s, mocksvc, ctx := initializeTest(t)

	input := make(map[string]string, 1)
	input["id1"] = "value1"

	successResp := buildSuccessResponse(http.StatusCreated, input)
	serverErr := getServerError()
	errResp := getErrResp()

	tests := []struct {
		desc     string
		key      string
		value    string
		mockResp *service.Response // mock response from get call
		mockErr  error
		err      error
	}{
		{"successful creation", "id1", "value1", successResp, nil, nil},
		{"request failure", "id1", "value1", errResp, nil, serverErr},
		{"Call to Post service Error", "id1", "value1", nil, serviceErr, serviceErr},
	}

	for i, tc := range tests {
		input := make(map[string]string)
		input[tc.key] = tc.value

		body, err := json.Marshal(&input)
		if err != nil {
			t.Error(err)

			continue
		}

		mocksvc.EXPECT().Post(ctx, "data", nil, body).Return(tc.mockResp, tc.mockErr)

		err = s.Set(ctx, tc.key, tc.value)

		assert.Equal(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

func TestService_CreateBindError(t *testing.T) {
	s, mockSvc, ctx := initializeTest(t)

	input := make(map[string]string, 1)
	input["id1"] = "value1"

	body, err := json.Marshal(&input)
	if err != nil {
		t.Errorf("Received unexpected error:\n%+v", err)
	}

	bindErr := getBindErr()
	mockResp := &service.Response{StatusCode: http.StatusBadRequest, Body: []byte(`invalid body`)}

	mockSvc.EXPECT().Post(gomock.Any(), "data", nil, body).Return(mockResp, nil)
	mockSvc.EXPECT().Bind(mockResp.Body, gomock.Any()).Return(errors.Error("bind error"))

	err = s.Set(ctx, "id1", "value1")

	assert.Equal(t, bindErr, err)
}

func Test_Get(t *testing.T) {
	s, mockSvc, ctx := initializeTest(t)

	successOutput := map[string]string{"id1": "value1"}

	successResp := buildSuccessResponse(http.StatusOK, successOutput)
	serverErr := getServerError()
	errResp := getErrResp()

	tests := []struct {
		desc     string
		mockResp *service.Response // mock response from get call
		mockErr  error             // mock error from get call
		resp     string
		err      error
	}{
		{"call to Get service throws error", nil, serviceErr, "", serviceErr},
		{"response with bad request", errResp, nil, "", serverErr},
		{"success case", successResp, nil, "value1", nil},
	}

	for i, tc := range tests {
		mockSvc.EXPECT().Get(gomock.Any(), "data/id1", nil).Return(tc.mockResp, tc.mockErr)

		resp, err := s.Get(ctx, "id1")

		assert.Equal(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)

		assert.Equal(t, tc.resp, resp, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

func TestService_GetBindError(t *testing.T) {
	s, mockSvc, ctx := initializeTest(t)

	mockResp := &service.Response{StatusCode: http.StatusBadRequest, Body: []byte(`invalid body`)}
	errResp := &service.Response{StatusCode: http.StatusOK, Body: []byte(`invalid body`)}

	bindErr := getBindErr()

	tests := []struct {
		desc     string
		mockResp *service.Response // mock response from get call
	}{
		{"error in binding error response", mockResp},
		{"error in binding data response", errResp},
	}

	for i, tc := range tests {
		mockSvc.EXPECT().Get(gomock.Any(), "data/id1", nil).Return(tc.mockResp, nil)
		mockSvc.EXPECT().Bind(tc.mockResp.Body, gomock.Any()).Return(errors.Error("bind error"))

		output, err := s.Get(ctx, "id1")

		assert.Equal(t, bindErr, err, "TEST[%d], failed.\n%s", i, tc.desc)

		assert.Equal(t, "", output, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

func TestKvData_Delete(t *testing.T) {
	s, mockSvc, ctx := initializeTest(t)

	mockResp := &service.Response{Body: nil, StatusCode: http.StatusNoContent}

	serverErr := getServerError()
	errResp := getErrResp()

	tests := []struct {
		desc     string
		key      string
		mockResp *service.Response
		mockErr  error
		err      error
	}{
		{"success", "id1", mockResp, nil, nil},
		{"failed", "id1", mockResp, serviceErr, serviceErr},
		{"response with bad request", "id1", errResp, nil, serverErr},
	}

	for i, tc := range tests {
		mockSvc.EXPECT().Delete(ctx, "data/"+tc.key, nil).Return(tc.mockResp, tc.mockErr)

		err := s.Delete(ctx, tc.key)

		assert.Equal(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

func TestService_DeleteBindError(t *testing.T) {
	s, mockSvc, ctx := initializeTest(t)

	bindErr := &errors.Response{
		Code:   "Bind Error",
		Reason: "failed to bind response",
		Detail: errors.Error("bind error"),
	}

	test := struct {
		desc     string
		mockResp *service.Response // mock response from get call
	}{"error in binding error response", &service.Response{StatusCode: http.StatusBadRequest,
		Body: []byte(`invalid body`)},
	}

	mockSvc.EXPECT().Delete(gomock.Any(), "data/id1", nil).Return(test.mockResp, nil)
	mockSvc.EXPECT().Bind(test.mockResp.Body, gomock.Any()).Return(errors.Error("bind error"))

	err := s.Delete(ctx, "id1")

	assert.Equal(t, bindErr, err, "TEST failed.\n%s", test.desc)
}

// buildSuccessResponse builds the success response
func buildSuccessResponse(statusCode int, input map[string]string) *service.Response {
	resp := struct {
		Data map[string]string `json:"data"`
	}{Data: input}

	body, _ := json.Marshal(&resp)

	return &service.Response{
		Body:       body,
		StatusCode: statusCode,
	}
}

// getServerError returns a serverError
func getServerError() error {
	return errors.MultipleErrors{
		StatusCode: http.StatusInternalServerError,
		Errors: []error{&errors.Response{
			Code:     "Internal Server Error",
			Reason:   "connection timed out",
			DateTime: errors.DateTime{Value: "2021-11-03T11:01:13.124Z", TimeZone: "IST"},
		}}}
}

// getErrResp returns a HttpService response with error message
func getErrResp() (errResp *service.Response) {
	errResponse := []errors.Response{{Code: "Internal Server Error", Reason: "connection timed out",
		DateTime: errors.DateTime{Value: "2021-11-03T11:01:13.124Z", TimeZone: "IST"}}}

	resp := struct {
		Errors []errors.Response `json:"errors"`
	}{Errors: errResponse}

	body, _ := json.Marshal(resp)

	return &service.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       body,
	}
}

// getBindErr returns bind error
func getBindErr() error {
	return &errors.Response{
		Code:   "Bind Error",
		Reason: "failed to bind response",
		Detail: errors.Error("bind error"),
	}
}

type mockHTTPService struct {
	*mocks.MockHTTPService // embed mock of HTTPService interface
}

// override the Bind method
func (m mockHTTPService) Bind(resp []byte, i interface{}) error {
	if err := json.Unmarshal(resp, i); err != nil {
		return m.MockHTTPService.Bind(resp, i)
	}

	return nil
}
