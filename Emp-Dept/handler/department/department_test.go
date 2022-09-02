package department

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"developer.zopsmart.com/go/gofr/Emp-Dept/entities"
	"developer.zopsmart.com/go/gofr/Emp-Dept/service"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
	"developer.zopsmart.com/go/gofr/pkg/gofr/responder"
)

func TestPost(t *testing.T) {
	app := gofr.New()

	testCases := []struct {
		desc           string
		input          []byte
		expectedOutput interface{}
		err            error
	}{
		{"Valid input", []byte(`{"id":1, "name":"HR","floorNo": 1}`), entities.Department{1, "HR", 1}, nil},
		{"Invalid input", []byte(`{"id":0,"name":"Tech","floorNo":2}`), nil, errors.InvalidParam{Param: []string{"invalid id"}}},
	}

	for i, tc := range testCases {
		var dept entities.Department

		r := httptest.NewRequest("POST", "/department", bytes.NewReader(tc.input))
		w := httptest.NewRecorder()

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(res, req, app)

		ctrl := gomock.NewController(t)
		mockIndex := service.NewMockDepartment(ctrl)

		err := ctx.Bind(&dept)
		if err != nil {
			log.Println(err)
		}

		mockIndex.EXPECT().Post(ctx, dept).Return(tc.expectedOutput, tc.err)

		mockService := Handler{service: mockIndex}

		actualOutput, _ := mockService.Post(ctx)
		if !reflect.DeepEqual(actualOutput, tc.expectedOutput) {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, actualOutput, tc.expectedOutput)
		}
	}
}

func TestPostWithEdgeCases(t *testing.T) {
	app := gofr.New()

	testCases := []struct {
		desc           string
		input          []byte
		expectedOutput interface{}
	}{
		{"for Unmarshal error", []byte(`{"id":"2","name":hr,"floorNo":2}`), nil},
	}

	for i, tc := range testCases {
		r := httptest.NewRequest("POST", "/department", bytes.NewReader(tc.input))
		w := httptest.NewRecorder()

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(res, req, app)

		ctrl := gomock.NewController(t)
		mockIndex := service.NewMockDepartment(ctrl)
		mockService := Handler{service: mockIndex}

		actualOutput, _ := mockService.Post(ctx)
		if !reflect.DeepEqual(actualOutput, tc.expectedOutput) {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, actualOutput, tc.expectedOutput)
		}
	}
}

func TestPut(t *testing.T) {
	app := gofr.New()

	testcases := []struct {
		desc           string
		id             int
		dataToUpdate   []byte
		expectedOutput interface{}
		err            error
	}{
		{"valid case", 1, []byte(`{"id":1, "name":"tech","floorNo": 2}`), entities.Department{1, "TECH", 2}, nil},
		{"Invalid ID", 4, []byte(`{"ID":2,"Name":"tech","FloorNo":2}`), nil, errors.InvalidParam{Param: []string{"invalid id"}}},
	}

	for i, tc := range testcases {
		var dept entities.Department

		path := fmt.Sprintf("/department/%v", tc.id)

		r := httptest.NewRequest("PUT", path, bytes.NewReader(tc.dataToUpdate))
		w := httptest.NewRecorder()

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(res, req, app)

		ctrl := gomock.NewController(t)
		mockIndex := service.NewMockDepartment(ctrl)
		mockService := Handler{service: mockIndex}

		err := ctx.Bind(&dept)
		if err != nil {
			log.Println(err)
		}

		mockIndex.EXPECT().Put(ctx, tc.id, dept).Return(tc.expectedOutput, tc.err)

		actualOutput, _ := mockService.Put(ctx)
		if !reflect.DeepEqual(actualOutput, tc.expectedOutput) {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, actualOutput, tc.expectedOutput)
		}
	}
}

func TestPutWithEdgesCases(t *testing.T) {
	app := gofr.New()

	testcases := []struct {
		desc           string
		id             string
		dataToUpdate   []byte
		expectedOutput interface{}
		err            error
	}{
		{"Testcase for Unmarshal Error", "2", []byte(nil), nil, errors.InvalidParam{Param: []string{"invalid"}}},
		{"unconvertable id ", "abc", []byte(nil), nil, errors.InvalidParam{Param: []string{"invalid"}}},
	}

	for i, tc := range testcases {
		path := fmt.Sprintf("/department/%v", tc.id)

		r := httptest.NewRequest("PUT", path, bytes.NewReader(tc.dataToUpdate))
		w := httptest.NewRecorder()

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(res, req, app)

		ctrl := gomock.NewController(t)
		mockIndex := service.NewMockDepartment(ctrl)
		mockService := New(mockIndex)

		actualOutput, _ := mockService.Put(ctx)
		if !reflect.DeepEqual(actualOutput, tc.expectedOutput) {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, actualOutput, tc.expectedOutput)
		}
	}
}

func TestDelete(t *testing.T) {
	app := gofr.New()
	testcases := []struct {
		desc           string
		input          int
		expectedOutput int
		err            error
	}{
		{"id in db", 1, 204, nil},
		{"id not in db", 4, 404, errors.DB{Err: fmt.Errorf("id not found")}},
	}

	for i, tc := range testcases {
		path := fmt.Sprintf("/department/%v", tc.input)

		r := httptest.NewRequest("DELETE", path, nil)
		w := httptest.NewRecorder()

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(res, req, app)

		ctrl := gomock.NewController(t)
		mockIndex := service.NewMockDepartment(ctrl)
		mockService := Handler{service: mockIndex}

		mockIndex.EXPECT().Delete(ctx, tc.input).Return(tc.expectedOutput, tc.err)

		actualOutput, err := mockService.Delete(ctx)
		if err != tc.err {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, actualOutput, tc.expectedOutput)
		}
	}
}

func TestDeleteWithEdgeCases(t *testing.T) {
	app := gofr.New()
	testcases := []struct {
		desc           string
		input          interface{}
		expectedOutput int
		err            error
	}{
		{"unconvertable type", 1.4, 400, errors.InvalidParam{Param: []string{"Invalid id"}}},
	}

	for i, tc := range testcases {
		r := httptest.NewRequest(http.MethodDelete, "/department/"+fmt.Sprintf("%v", tc.input), nil)
		req := request.NewHTTPRequest(r)
		ctx := gofr.NewContext(nil, req, app)

		ctrl := gomock.NewController(t)
		mockIndex := service.NewMockDepartment(ctrl)
		mockService := Handler{service: mockIndex}

		_, err := mockService.Delete(ctx)
		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, err, tc.err)
		}
	}
}
