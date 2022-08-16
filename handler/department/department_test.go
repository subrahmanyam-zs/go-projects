package department

import (
	"EmployeeDepartment/entities"
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestEmployeePost(t *testing.T) {

	testCases := []struct {
		desc           string
		input          []byte
		expectedOutput []byte
	}{
		{"Valid input", []byte(`{"id":1, "name":"HR","floorNo": 1}`), []byte(`{"Id":1,"Name":"HR","FloorNo":1}`)},
		{"Invalid input", []byte(`{"id":0,"name":"Tech","floorNo":2}`), []byte("Invalid id")},
		{"for Unmarshal error", []byte(`{"id":"2","name":"hr","floorNo":2}`), []byte("Unmarshal Error")},
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

func (m mockService) Create(department entities.Department) (entities.Department, error) {

	if department.Id == 0 {
		return entities.Department{}, errors.New("error")
	}
	return entities.Department{1, "HR", 1}, nil
}

func TestPutHandler(t *testing.T) {
	testcases := []struct {
		desc           string
		input          int
		dataToUpdate   []byte
		expectedOutput []byte
	}{
		{"valid case", 1, []byte(`{"id":1, "name":"tech","floorNo": 2}`), []byte(`{"Id":1,"Name":"tech","FloorNo":2}`)},
		{"Testcase for Unmarshal Error", 2, []byte(nil), []byte("Unmarshal Error")},
		{"Invalid Id", 4, []byte(`{"Id":2,"Name":"tech","FloorNo":2}`), []byte("Id not found")},
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

func (m mockService) Update(id int, department entities.Department) (entities.Department, error) {
	if id == 1 {
		return entities.Department{1, "tech", 2}, nil
	}
	return entities.Department{}, errors.New("error")
}

func TestDeleteHandler(t *testing.T) {
	testcases := []struct {
		desc           string
		input          int
		expectedOutput int
	}{
		{"Valid Id", 1, 204},
		{"Id not found", 4, 404},
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

func (m mockService) Delete(id int) (int, error) {
	if id == 1 {
		return http.StatusNoContent, nil
	}
	return http.StatusNotFound, errors.New("error")

}

type mockService struct {
	id int
}
