package employee

import (
	"bytes"
	"fmt"
	"log"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	uuidMongo "go.mongodb.org/mongo-driver/x/mongo/driver/uuid"

	"developer.zopsmart.com/go/gofr/Emp-Dept/entities"
	"developer.zopsmart.com/go/gofr/Emp-Dept/service"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
	"developer.zopsmart.com/go/gofr/pkg/gofr/responder"
)

func TestPost(t *testing.T) {
	testcases := []struct {
		desc           string
		input          []byte
		expectedOutput interface{}
		err            error
	}{
		{"Valid input", []byte(`{"Name": "jason", "Dob": "12-06-2002","City": "Bangalore","Major":"CSE","DId": 1}`),
			[]byte(fmt.Sprintf(`{"ID":"%v","Name":"%v","Dob":"%v","City":"%v","Majors":"%v","DId":%v}`,
				"123e4567-e89b-12d3-a456-426614174000", "jason", "12-06-2002", "Bangalore", "CSE", 1)), nil},
		{"Empty name", []byte(`{"Name": "","Dob": "12-06-2002","City": "Bangalore","Major":"CSE","DepId": 1}`),
			nil, errors.InvalidParam{Param: []string{"invalid details"}}},
	}

	for i, tc := range testcases {
		var emp entities.Employee

		r := httptest.NewRequest("POST", "/employee", bytes.NewReader(tc.input))
		w := httptest.NewRecorder()

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(res, req, gofr.New())

		ctrl := gomock.NewController(t)
		mockIndex := service.NewMockEmployee(ctrl)

		err := ctx.Bind(&emp)
		if err != nil {
			log.Println(err)
		}

		mockIndex.EXPECT().Post(ctx, emp).Return(tc.expectedOutput, tc.err)

		mockService := Handler{service: mockIndex}
		actualOutput, err := mockService.Post(ctx)

		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("testcase %d failed got %v \n expected %v", i+1, actualOutput, tc.expectedOutput)
		}
	}
}

func TestPostWithEdgeCases(t *testing.T) {
	testcases := []struct {
		desc           string
		input          []byte
		expectedOutput interface{}
		err            error
	}{
		{"for Unmarshal error", []byte(nil),
			[]byte("Invalid body"), errors.InvalidParam{Param: []string{"invalid parameters"}}},
	}

	for i, tc := range testcases {
		r := httptest.NewRequest("POST", "/employee/", bytes.NewReader(tc.input))
		w := httptest.NewRecorder()

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(res, req, gofr.New())

		ctrl := gomock.NewController(t)
		mockIndex := service.NewMockEmployee(ctrl)

		mockService := Handler{service: mockIndex}
		_, err := mockService.Post(ctx)

		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("testcase %d failed got %v \n expected %v", i+1, err, tc.err)
		}
	}
}

func TestPut(t *testing.T) {
	testcases := []struct {
		desc           string
		input          string
		dataToUpdate   []byte
		expectedOutput interface{}
		err            error
	}{
		{"Valid case", "123e4567-e89b-12d3-a456-426614174000",
			[]byte(fmt.Sprintf(`{"Name":"jason","Dob":"12-06-2002","City":"Bangalore","Majors":"CSE","DId":2}`)),
			entities.Employee{uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), "jason", "12-06-2002", "Bangalore", "CSE", 2}, nil},
		{"InValid case", "123e4567-e89b-12d3-a456-426614174000",
			[]byte(fmt.Sprintf(`{"Name":"jason","Dob":"12-06-2004","City":"Bangalore","Majors":"CSE","DId":2}`)),
			entities.Employee{uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), "jason", "12-06-2002", "Bangalore", "CSE", 2},
			errors.InvalidParam{Param: []string{"invalid details"}}},
	}

	for i, tc := range testcases {
		var emp entities.Employee

		r := httptest.NewRequest("PUT", "/employee/", bytes.NewReader(tc.dataToUpdate))
		w := httptest.NewRecorder()

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(res, req, gofr.New())
		ctx.SetPathParams(map[string]string{
			"id": tc.input,
		})

		ctrl := gomock.NewController(t)
		mockIndex := service.NewMockEmployee(ctrl)

		err := ctx.Bind(&emp)
		if err != nil {
			log.Println(err)
		}

		mockIndex.EXPECT().Put(ctx, uuid.MustParse(tc.input), emp).Return(tc.expectedOutput, tc.err)

		mockService := New(mockIndex)

		_, err = mockService.Put(ctx)
		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("testcase %d failed got %v \n expected %v", i+1, err, tc.err)
		}
	}
}

func TestPutWithEdgeCases(t *testing.T) {
	testcases := []struct {
		desc           string
		input          string
		dataToUpdate   []byte
		expectedOutput interface{}
		err            error
	}{
		{"invalid uuid format", "123e4567-e89b-12d3a456-426614174000",
			[]byte(fmt.Sprintf(`{"Name":"jason","Dob":"12-06-2002","City":"Bangalore","Majors":"CSE","DId":1}`)),
			nil, errors.InvalidParam{Param: []string{"invalid id"}}},
		{"unmarshal error", "123e4567-e89b-12d3-a456-426614174000",
			[]byte(fmt.Sprintf(`{"Name":jason,"Dob":"12-06-2002","City":"Bangalore","Majors":"CSE","DId":1}`)),
			nil, errors.InvalidParam{Param: []string{"invalid details"}}},
	}

	for i, tc := range testcases {
		r := httptest.NewRequest("PUT", "/employee/", bytes.NewReader(tc.dataToUpdate))
		w := httptest.NewRecorder()

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(res, req, gofr.New())
		ctx.SetPathParams(map[string]string{
			"id": tc.input,
		})

		ctrl := gomock.NewController(t)
		mockIndex := service.NewMockEmployee(ctrl)

		mockService := Handler{service: mockIndex}

		_, err := mockService.Put(ctx)
		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("testcase %d failed got %v \n expected %v", i+1, err, tc.err)
		}
	}
}

func TestDeleteHandler(t *testing.T) {
	app := gofr.New()
	testcases := []struct {
		desc           string
		input          string
		expectedOutput int
		err            error
	}{
		{"Valid ID", "123e4567-e89b-12d3-a456-426614174000", 204, nil},
		{"ID not found", "123e4567-e89b-12d3-a456-426614174020", 404, errors.DB{Err: fmt.Errorf("err")}},
	}

	for i, tc := range testcases {
		r := httptest.NewRequest("DELETE", "/employee/", nil)
		w := httptest.NewRecorder()

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(res, req, app)
		ctx.SetPathParams(map[string]string{
			"id": tc.input,
		})

		ctrl := gomock.NewController(t)
		mockIndex := service.NewMockEmployee(ctrl)

		mockIndex.EXPECT().Delete(ctx, uuid.MustParse(tc.input)).Return(tc.expectedOutput, tc.err)

		mockService := Handler{service: mockIndex}

		_, err := mockService.Delete(ctx)
		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("testcase %d failed got %v \n expected %v", i+1, err, tc.err)
		}
	}
}

func TestDeleteWithEdgeCases(t *testing.T) {
	app := gofr.New()
	testcases := []struct {
		desc           string
		input          string
		expectedOutput int
		err            error
	}{
		{"invalid uuid format", "123e4567e89b-12d3-a456426614174020", 400, errors.InvalidParam{Param: []string{"invalid id"}}},
	}

	for i, tc := range testcases {
		r := httptest.NewRequest("DELETE", "/employee/", nil)
		w := httptest.NewRecorder()

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(res, req, app)
		ctx.SetPathParams(map[string]string{
			"id": tc.input,
		})

		ctrl := gomock.NewController(t)
		mockIndex := service.NewMockEmployee(ctrl)

		mockService := Handler{service: mockIndex}

		actualOutput, err := mockService.Delete(ctx)
		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("testcase %d failed got %v \n expected %v", i+1, actualOutput, tc.expectedOutput)
		}
	}
}

func TestGet(t *testing.T) {
	testcases := []struct {
		desc           string
		input          string
		expectedOutput interface{}
		err            error
	}{
		{"Valid Input", "123e4567-e89b-12d3-a456-426614174000",
			entities.EmpDept{uuid.UUID(uuidMongo.UUID(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))),
				"jason", "12-06-1998", "Kochi", "CSE",
				entities.Department{2, "TECH", 2}}, nil},
		{"InValid Input", "123e4567-e89b-12d3-a456-426614174000", nil, errors.DB{Err: fmt.Errorf("id not found")}},
	}
	for i, tc := range testcases {
		r := httptest.NewRequest("GET", "/employee/", nil)
		w := httptest.NewRecorder()

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(res, req, gofr.New())
		ctx.SetPathParams(map[string]string{
			"id": tc.input,
		})

		ctrl := gomock.NewController(t)
		mockIndex := service.NewMockEmployee(ctrl)

		mockIndex.EXPECT().Get(ctx, uuid.MustParse(tc.input)).Return(tc.expectedOutput, tc.err)

		mockService := Handler{service: mockIndex}

		_, err := mockService.Get(ctx)
		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("testcase %d failed got %v \n expected %v", i+1, err, tc.err)
		}
	}
}

func TestGetWithEdgeCases(t *testing.T) {
	app := gofr.New()
	testcases := []struct {
		desc           string
		input          string
		expectedOutput interface{}
		err            error
	}{
		{"invalid uuid format", "123e4567e89b-12d3-a456426614174020", nil, errors.InvalidParam{Param: []string{"invalid id"}}},
	}

	for i, tc := range testcases {
		r := httptest.NewRequest("GET", fmt.Sprintf("/employee/%v", tc.input), nil)
		w := httptest.NewRecorder()

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(res, req, app)
		ctx.SetPathParams(map[string]string{
			"id": tc.input,
		})

		ctrl := gomock.NewController(t)
		mockIndex := service.NewMockEmployee(ctrl)

		mockService := Handler{service: mockIndex}

		actualOutput, err := mockService.Get(ctx)
		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("testcase %d failed got %v \n expected %v", i+1, actualOutput, tc.expectedOutput)
		}
	}
}

func TestGetAll(t *testing.T) {
	testcases := []struct {
		desc              string
		name              string
		includeDepartment bool
		expectedOutput    interface{}
		err               error
	}{
		{"valid name and including dept ", "jason", true, entities.EmpDept{
			uuid.UUID(uuidMongo.UUID(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))), "jason", "12-06-1998", "Kochi", "CSE",
			entities.Department{2, "TECH", 2}}, nil},
		{"invalid name ", "", true, nil, errors.InvalidParam{Param: []string{"invalid name"}}},
	}
	for i, tc := range testcases {
		r := httptest.NewRequest("GET", fmt.Sprintf("/employee?name=%v&includeDepartment=%v", tc.name, tc.includeDepartment), nil)
		w := httptest.NewRecorder()

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(res, req, gofr.New())

		ctrl := gomock.NewController(t)
		mockIndex := service.NewMockEmployee(ctrl)

		mockIndex.EXPECT().GetAll(ctx, tc.name, tc.includeDepartment).Return(tc.expectedOutput, tc.err)

		mockService := Handler{service: mockIndex}

		_, err := mockService.GetAll(ctx)
		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("testcase %d failed got %v \n expected %v", i+1, err, tc.err)
		}
	}
}

func TestGetAllWithEdgeCases(t *testing.T) {
	testcases := []struct {
		desc              string
		name              string
		includeDepartment string
		expectedOutput    interface{}
		err               error
	}{
		{"parsebool error", "jason", "abcd",
			nil, errors.InvalidParam{Param: []string{"unconvertable type to bool"}}},
	}
	for i, tc := range testcases {
		r := httptest.NewRequest("GET", fmt.Sprintf("/employee?name=%v&includeDepartment=%v", tc.name, tc.includeDepartment), nil)
		w := httptest.NewRecorder()

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(res, req, gofr.New())

		ctrl := gomock.NewController(t)
		mockIndex := service.NewMockEmployee(ctrl)

		mockService := Handler{service: mockIndex}

		_, err := mockService.GetAll(ctx)
		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("testcase %d failed got %v \n expected %v", i+1, err, tc.err)
		}
	}
}
