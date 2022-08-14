package employee

import (
	entities2 "EmployeeDepartment/entities"
	"bytes"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var uid uuid.UUID = uuid.New()

func TestPostHandler(t *testing.T) {
	testcases := []struct {
		desc           string
		input          []byte
		expectedOutput []byte
	}{
		{"Valid input", []byte(`{"Name": "jason", "Dob": "12-06-2002","City": "Bangalore","Major":"CSE","DId": 1}`),
			[]byte(fmt.Sprintf(`{"Id":"%v","Name":"%v","Dob":"%v","City":"%v","Majors":"%v","DId":%v}`, uid, "jason", "12-06-2002", "Bangalore", "CSE", 1))},
		{"Empty name", []byte(`{"Name": "","Dob": "12-06-2002","City": "Bangalore","Major":"CSE","DepId": 1}`), []byte(`Invalid Body`)},
		{"for Unmarshal error", []byte(`{"Name":jason,"Dob":"12-06-2002","City":"Bangalore","Major":"CSE","DepId":1}`), []byte("Unmarshal Error")},
	}

	for i, tc := range testcases {

		reader := bytes.NewReader(tc.input)
		req := httptest.NewRequest(http.MethodPost, "/post", reader)
		res := httptest.NewRecorder()

		a := New(mockService{})

		a.PostHandler(res, req)

		if !reflect.DeepEqual(res.Body, bytes.NewBuffer(tc.expectedOutput)) {
			t.Errorf("testcase %d failed got %v \n expected %v", i+1, res.Body.String(), string(tc.expectedOutput))
		}
	}
}

func (m mockService) Create(employee entities2.Employee) (entities2.Employee, error) {
	name := employee.Name
	if name == "" {
		return entities2.Employee{}, errors.New("error")
	}
	return entities2.Employee{uid, "jason", "12-06-2002", "Bangalore", "CSE", 1}, nil
}

func TestGetHandler(t *testing.T) {
	testcases := []struct {
		desc           string
		input          string
		expectedOutput []byte
	}{
		{"Valid Input", uid.String(), []byte(fmt.Sprintf(`{"Id":"%v","Name":"jason","Dob":"12-06-2002","City":"Bangalore","Majors":"CSE","DId":1}`, uid))},
		{"Invalid Id", "00000-12223-122-1222323-2133", []byte("Invalid Id")},
		{"Id not found", "123e4567-e89b-12d3-a456-426614174000", []byte("Id not found")},
	}

	for i, tc := range testcases {
		path := fmt.Sprintf("/employee/%v", tc.input)
		req := httptest.NewRequest("GET", path, nil)
		res := httptest.NewRecorder()

		a := New(mockService{})

		a.GetHandler(res, req)

		if !reflect.DeepEqual(res.Body, bytes.NewBuffer(tc.expectedOutput)) {
			t.Errorf("testcase %d failed got %v \n expected %v", i+1, res.Body.String(), string(tc.expectedOutput))
		}
	}
}

func (m mockService) Read(id uuid.UUID) (entities2.Employee, error) {
	if uid == id {
		return entities2.Employee{uid, "jason", "12-06-2002", "Bangalore", "CSE", 1}, nil
	}
	return entities2.Employee{}, errors.New("error")
}

func TestPutHandler(t *testing.T) {
	testcases := []struct {
		desc           string
		input          string
		dataToUpdate   []byte
		expectedOutput []byte
	}{
		{"Valid case", uid.String(), []byte(fmt.Sprintf(`{"Id":"%v","Name":"jason","Dob":"12-06-2002","City":"Bangalore","Majors":"CSE","DId":1}`, uid)),
			[]byte(fmt.Sprintf(`{"Id":"%v","Name":"jason","Dob":"12-06-2002","City":"Bangalore","Majors":"CSE","DId":1}`, uid))},
		{"for unmarshal error", uid.String(), []byte(nil), []byte("Unmarshall error")},
		{"Invalid Id", "123e4567-e89b-12d3-a456-426614174000", []byte(fmt.Sprintf(`{"Id":"%v","Name":"jason","Dob":"12-06-2002","City":"Bangalore","Majors":"CSE","DId":1}`, uid)),
			[]byte(`Id not found`)},
	}

	for i, tc := range testcases {
		path := fmt.Sprintf("/employee/%v", tc.input)
		req := httptest.NewRequest(http.MethodPut, path, bytes.NewReader(tc.dataToUpdate))
		res := httptest.NewRecorder()

		a := New(mockService{})

		a.PutHandler(res, req)
		if !reflect.DeepEqual(res.Body, bytes.NewBuffer(tc.expectedOutput)) {
			t.Errorf("testcase %d failed got %v \n expected %v", i+1, res.Body.String(), string(tc.expectedOutput))
		}
	}
}

func (m mockService) Update(id uuid.UUID, employee entities2.Employee) (entities2.Employee, error) {
	if uid == id {
		return entities2.Employee{uid, "jason", "12-06-2002", "Bangalore", "CSE", 1}, nil
	}
	return entities2.Employee{}, errors.New("error")
}

func TestDeleteHandler(t *testing.T) {
	testcases := []struct {
		desc           string
		input          string
		expectedOutput int
	}{
		{"Valid Id", uid.String(), 204},
		{"Id not found", "123e4567-e89b-12d3-a456-426614174000", 404},
	}

	for i, tc := range testcases {
		path := fmt.Sprintf("/employee/%v", tc.input)
		req := httptest.NewRequest(http.MethodDelete, path, nil)
		res := httptest.NewRecorder()

		a := New(mockService{})

		a.DeleteHandler(res, req)

		if res.Result().StatusCode != tc.expectedOutput {
			t.Errorf("testcase %d failed got %v \n expected %v", i+1, res.Body.String(), tc.expectedOutput)
		}
	}
}

func (m mockService) Delete(id uuid.UUID) (int, error) {
	if id == uid {
		return http.StatusNoContent, nil
	}
	return http.StatusNotFound, errors.New("error")
}

func TestGetAll(t *testing.T) {
	testcases := []struct {
		desc              string
		name              string
		includeDepartment bool
		expectedOutput    []byte
	}{
		{"valid name and including dept ", "jason", true,
			[]byte(fmt.Sprintf(`{"Id":"%v","Name":"jason","Dob":"12-06-1998","City":"Bangalore","Majors":"CSE","Dept":{"Id":1,"Name":"HR","FloorNo":1}}`, uid))},
		{"valid name and not including dept ", "jason", false,
			[]byte(fmt.Sprintf(`{"Id":"%v","Name":"jason","Dob":"12-06-1998","City":"Bangalore","Majors":"CSE","Dept":{"Id":0,"Name":"","FloorNo":0}}`, uid))},
		{"invalid name and includeDepartment is true", "", true,
			[]byte("Unmarshal Error")},
	}
	for i, tc := range testcases {
		path := fmt.Sprintf("/employee?name=%v&includeDepartment=%v", tc.name, tc.includeDepartment)
		req := httptest.NewRequest(http.MethodGet, path, nil)
		res := httptest.NewRecorder()

		a := New(mockService{})

		a.GetAll(res, req)

		if !reflect.DeepEqual(res.Body, bytes.NewBuffer(tc.expectedOutput)) {
			t.Errorf("testcase %d failed got %v \n expected %v", i+1, res.Body.String(), string(tc.expectedOutput))
		}
	}
}

func (m mockService) ReadAll(name string, includeDepartment bool) (entities2.EmployeeAndDepartment, error) {
	if (name != "") && (includeDepartment == true) {
		return entities2.EmployeeAndDepartment{uid.String(), "jason", "12-06-1998", "Bangalore", "CSE", entities2.Department{1, "HR", 1}}, nil
	} else if (name != "") && (includeDepartment == false) {
		return entities2.EmployeeAndDepartment{uid.String(), "jason", "12-06-1998", "Bangalore", "CSE", entities2.Department{}}, nil
	}
	return entities2.EmployeeAndDepartment{}, errors.New("error")
}

type mockService struct{}
