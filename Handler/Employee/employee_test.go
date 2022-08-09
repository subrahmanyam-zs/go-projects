package Employee

import (
	"EmployeeDepartment/Handler/Entities"
	"bytes"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var uid uuid.UUID = generateUUID()

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

		a := New(mockDatastore{})

		a.PostHandler(res, req)

		if !reflect.DeepEqual(res.Body, bytes.NewBuffer(tc.expectedOutput)) {
			t.Errorf("testcase %d failed got %v \n expected %v", i+1, res.Body.String(), string(tc.expectedOutput))
		}
	}
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

		a := New(mockDatastore{})

		a.GetHandler(res, req)

		if !reflect.DeepEqual(res.Body, bytes.NewBuffer(tc.expectedOutput)) {
			t.Errorf("testcase %d failed got %v \n expected %v", i+1, res.Body.String(), string(tc.expectedOutput))
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
		{"Valid case", uid.String(), []byte(fmt.Sprintf(`{"Id":"%v","Name":"jason","Dob":"12-06-2002","City":"Bangalore","Majors":"CSE","DId":1}`, uid)), []byte(fmt.Sprintf(`{"Id":"%v","Name":"jason","Dob":"12-06-2002","City":"Bangalore","Majors":"CSE","DId":1}`, uid))},
		{"for unmarshal error", uid.String(), []byte(nil), []byte("Unmarshall error")},
		//{"Invalid Id", "123e4567-e89b-12d3-a456-426614174000", []byte(fmt.Sprintf(`{"Id":"%v","Name":"jason","Dob":"12-06-2002","City":"Bangalore","Majors":"CSE","DId":1}`, uid)), []byte(`{"Id":"00000000-0000-0000-0000-000000000000","Name":"","Dob":"","City":"","Majors":"","DId":0} `)},
	}

	for i, tc := range testcases {
		path := fmt.Sprintf("/employee/%v", tc.input)
		req := httptest.NewRequest(http.MethodPut, path, bytes.NewReader(tc.dataToUpdate))
		res := httptest.NewRecorder()

		a := New(mockDatastore{})

		a.PutHandler(res, req)
		if !reflect.DeepEqual(res.Body, bytes.NewBuffer(tc.expectedOutput)) {
			t.Errorf("testcase %d failed got %v \n expected %v", i+1, res.Body.String(), string(tc.expectedOutput))
		}
	}
}

type mockDatastore struct{}

func (m mockDatastore) Create(employee Entities.Employee) (Entities.Employee, error) {
	name := employee.Name
	if name == "" {
		return Entities.Employee{}, errors.New("error")
	}
	return Entities.Employee{uid, "jason", "12-06-2002", "Bangalore", "CSE", 1}, nil
}

func (m mockDatastore) Read(id uuid.UUID) (Entities.Employee, error) {
	if uid == id {
		return Entities.Employee{uid, "jason", "12-06-2002", "Bangalore", "CSE", 1}, nil
	}
	return Entities.Employee{}, errors.New("error")
}

func (m mockDatastore) Update(id uuid.UUID, employee Entities.Employee) (Entities.Employee, error) {
	if uid == id {
		return Entities.Employee{uid, "jason", "12-06-2002", "Bangalore", "CSE", 1}, nil
	}
	return Entities.Employee{}, errors.New("error")
}
