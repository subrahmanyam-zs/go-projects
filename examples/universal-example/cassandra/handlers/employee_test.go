package handlers

import (
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"developer.zopsmart.com/go/gofr/examples/universal-example/cassandra/entity"
	"developer.zopsmart.com/go/gofr/examples/universal-example/cassandra/store"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
)

func initializeHandlerTest(t *testing.T) (*store.MockEmployee, employee, *gofr.Gofr) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	employeeStore := store.NewMockEmployee(ctrl)
	employee := New(employeeStore)
	k := gofr.New()

	return employeeStore, employee, k
}

func TestCassandraEmployee_Get(t *testing.T) {
	tests := []struct {
		queryParams  string
		expectedResp []entity.Employee
		mockErr      error
	}{
		{"id=1", []entity.Employee{{ID: 1, Name: "Rohan", Phone: "01222", Email: "rohan@zopsmart.com", City: "Berlin"}}, nil},
		{"id=1&name=Rohan&phone=01222&email=rohan@zopsmart.com&city=Berlin",
			[]entity.Employee{{ID: 1, Name: "Rohan", Phone: "01222", Email: "rohan@zopsmart.com", City: "Berlin"}}, nil},
		{"", []entity.Employee{{ID: 1, Name: "Rohan", Phone: "01222", Email: "rohan@zopsmart.com", City: "Berlin"},
			{ID: 2, Name: "Aman", Phone: "22234", Email: "aman@zopsmart.com", City: "florida"}}, nil},
		{"id=7&name=Sunita", nil, nil},
	}

	employeeStore, employee, k := initializeHandlerTest(t)

	for i, tc := range tests {
		r := httptest.NewRequest("GET", "/employees?"+tc.queryParams, nil)
		req := request.NewHTTPRequest(r)
		context := gofr.NewContext(nil, req, k)

		employeeStore.EXPECT().Get(gomock.Any(), gomock.Any()).Return(tc.expectedResp)

		resp, err := employee.Get(context)
		assert.Equal(t, tc.mockErr, err, i)
		assert.Equal(t, tc.expectedResp, resp, i)
	}
}

func TestCassandraEmployee_Create(t *testing.T) {
	tests := []struct {
		query        string
		expectedResp interface{}
		mockErr      error
	}{
		{`{"id": 3, "name":"Shasank", "phone": "01567", "email":"shasank@zopsmart.com", "city":"Banglore"}`,
			[]entity.Employee{{ID: 3, Name: "Shasank", Phone: "01567", Email: "shasank@zopsmart.com", City: "Banglore"}}, nil},
		{`{"id": 4, "name":"Jay", "phone": "01933", "email":"jay@mahindra.com", "city":"Gujrat"}`,
			[]entity.Employee{{ID: 4, Name: "Jay", Phone: "01933", Email: "jay@mahindra.com", City: "Gujrat"}}, nil},
	}

	employeeStore, employee, k := initializeHandlerTest(t)

	for i, tc := range tests {
		input := strings.NewReader(tc.query)
		r := httptest.NewRequest("POST", "/dummy", input)
		req := request.NewHTTPRequest(r)
		context := gofr.NewContext(nil, req, k)

		employeeStore.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil)
		employeeStore.EXPECT().Create(gomock.Any(), gomock.Any()).Return(tc.expectedResp, tc.mockErr)

		resp, err := employee.Create(context)
		assert.Equal(t, tc.mockErr, err, i)
		assert.Equal(t, tc.expectedResp, resp, i)
	}
}

func TestCassandraEmployee_Create_InvalidInput_JsonError(t *testing.T) {
	tests := []struct {
		callGet       bool
		query         string
		expectedResp  interface{}
		mockGetOutput []entity.Employee
		mockErr       error
	}{
		// Invalid Input
		{true, `{"id": 2, "name": "Aman", "phone": "22234", "email": "aman@zopsmart.com", "city": "Florida"}`,
			nil, []entity.Employee{{ID: 2, Name: "Aman", Phone: "22234", Email: "aman@zopsmart.com", City: "Florida"}},
			errors.EntityAlreadyExists{}},
		// JSON Error
		{false, `{"id":    "2", "name":   "Aman", "phone": "22234", "email": "aman@zopsmart.com", "city": "Florida"}`, nil, nil,
			&json.UnmarshalTypeError{Value: "string", Type: reflect.TypeOf(2), Offset: 13, Struct: "Employee", Field: "id"}},
	}

	employeeStore, employee, k := initializeHandlerTest(t)

	for i, tc := range tests {
		input := strings.NewReader(tc.query)
		r := httptest.NewRequest("POST", "/dummy", input)
		req := request.NewHTTPRequest(r)
		context := gofr.NewContext(nil, req, k)

		if tc.callGet {
			employeeStore.EXPECT().Get(gomock.Any(), gomock.Any()).Return(tc.mockGetOutput).AnyTimes()
		}

		resp, err := employee.Create(context)
		assert.Equal(t, tc.mockErr, err, i)
		assert.Equal(t, tc.expectedResp, resp, i)
	}
}
