package department

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"EmployeeDepartment/entities"
	"EmployeeDepartment/errorsHandler"
)

type mockService struct{}

func TestEmployeePost(t *testing.T) {
	testCases := []struct {
		desc           string
		input          []byte
		expectedOutput []byte
	}{
		{"Valid input", []byte(`{"id":1, "name":"HR","floorNo": 1}`), []byte(`{"ID":1,"Name":"HR","FloorNo":1}`)},
		{"Invalid input", []byte(`{"id":0,"name":"Tech","floorNo":2}`), []byte("Invalid id")},
		{"for Unmarshal error", []byte(`{"id":"2","name":"hr","floorNo":2}`), []byte("Invalid body")},
	}

	for i, tc := range testCases {
		req := httptest.NewRequest("POST", "/employee", bytes.NewReader(tc.input))
		w := httptest.NewRecorder()

		a := New(mockService{})

		a.PostHandler(w, req)

		if !reflect.DeepEqual(w.Body, bytes.NewBuffer(tc.expectedOutput)) {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, w.Body.String(), string(tc.expectedOutput))
		}
	}
}

func TestPutHandler(t *testing.T) {
	testcases := []struct {
		desc           string
		input          string
		dataToUpdate   []byte
		expectedOutput []byte
	}{
		{"valid case", "1", []byte(`{"id":1, "name":"tech","floorNo": 2}`), []byte(`{"ID":1,"Name":"tech","FloorNo":2}`)},
		{"Testcase for Unmarshal Error", "2", []byte(nil), []byte("Invalid body")},
		{"Invalid ID", "4", []byte(`{"ID":2,"Name":"tech","FloorNo":2}`), []byte("id not found")},
		{"unconvertable string to int", "one", []byte(`{"ID":2,"Name":"tech","FloorNo":2}`), []byte("Invalid id")},
		{"invalid case", "2", []byte(`{"id":1, "name":"tech","floorNo": 2}`), []byte(`ID not found`)},
	}

	for i, tc := range testcases {
		path := fmt.Sprintf("/department/%v", tc.input)
		req := httptest.NewRequest(http.MethodPut, path, bytes.NewReader(tc.dataToUpdate))
		res := httptest.NewRecorder()

		a := New(mockService{})

		a.PutHandler(res, req)

		if !reflect.DeepEqual(res.Body, bytes.NewBuffer(tc.expectedOutput)) {
			t.Errorf("[TEST%d]Failed. Got %v\nExpected %v\n", i+1, res.Body.String(), string(tc.expectedOutput))
		}
	}
}

func TestDeleteHandler(t *testing.T) {
	testcases := []struct {
		desc           string
		input          string
		expectedOutput int
	}{
		{"Valid ID", "1", 204},
		{"ID not found", "4", 404},
		{"unconvertable string to type int", "one", 400},
	}
	for i, tc := range testcases {
		path := fmt.Sprintf("/department/%v", tc.input)
		req := httptest.NewRequest(http.MethodDelete, path, nil)
		res := httptest.NewRecorder()

		a := New(mockService{})
		a.DeleteHandler(res, req)

		if res.Result().StatusCode != tc.expectedOutput {
			t.Errorf("[TEST%d]Failed. Got %v\nExpected %v\n", i+1, res.Result().StatusCode, tc.expectedOutput)
		}
	}
}

func (m mockService) Create(ctx context.Context, department entities.Department) (entities.Department, error) {
	if department.ID == 0 {
		return entities.Department{}, errors.New("Invalid id")
	}

	return entities.Department{ID: 1, Name: "HR", FloorNo: 1}, nil
}

func (m mockService) Update(ctx context.Context, id int, department entities.Department) (entities.Department, error) {
	if id == 1 {
		return entities.Department{ID: 1, Name: "tech", FloorNo: 2}, nil
	}

	return entities.Department{}, errors.New("ID not found")
}

func (m mockService) Delete(ctx context.Context, id int) (int, error) {
	if id == 1 {
		return http.StatusNoContent, nil
	}

	return http.StatusNotFound, &errorsHandler.IDNotFound{Msg: "ID not found"}
}

func (m mockService) GetDepartment(ctx context.Context, id int) (entities.Department, error) {
	if id == 1 || id == 2 || id == 3 {
		return entities.Department{}, nil
	}
	return entities.Department{}, &errorsHandler.IDNotFound{Msg: "id not found"}
}
