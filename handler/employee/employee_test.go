package employee

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/google/uuid"

	entities "EmployeeDepartment/entities"
	"EmployeeDepartment/errorsHandler"
	"EmployeeDepartment/store"
)

type mockService struct{}

func TestPostHandler(t *testing.T) {
	testcases := []struct {
		desc           string
		input          []byte
		expectedOutput []byte
	}{
		{"Valid input", []byte(`{"Name": "jason", "Dob": "12-06-2002","City": "Bangalore","Major":"CSE","DId": 1}`),
			[]byte(fmt.Sprintf(`{"ID":"%v","Name":"%v","Dob":"%v","City":"%v","Majors":"%v","DId":%v}`,
				"123e4567-e89b-12d3-a456-426614174000", "jason", "12-06-2002", "Bangalore", "CSE", 1))},
		{"Empty name", []byte(`{"Name": "","Dob": "12-06-2002","City": "Bangalore","Major":"CSE","DepId": 1}`), []byte(`error`)},
		{"for Unmarshal error", []byte(`{"Name":jason,"Dob":"12-06-2002","City":"Bangalore","Major":"CSE","DepId":1}`), []byte("Invalid body")},
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

func TestGetHandler(t *testing.T) {
	testcases := []struct {
		desc           string
		input          string
		expectedOutput []byte
	}{
		{"Valid Input", "123e4567-e89b-12d3-a456-426614174000", []byte(fmt.Sprintf(
			`{"ID":"%v","Name":"jason","Dob":"12-06-2002","City":"Bangalore","Majors":"CSE","Dept":{"ID":2,"Name":"CSE","FloorNo":2}}`,
			"123e4567-e89b-12d3-a456-426614174000"))},
		{"Invalid ID", "00000-12223-122-1222323-2133", []byte("Invalid ID")},
		{"ID not found", "123e4567-e89b-12d3-a456-426614174200", []byte("error")},
		{"Invalid ID", "123e4567-e89b-12d3-a4565426614174000", []byte("Invalid id")},
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

func TestPutHandler(t *testing.T) {
	testcases := []struct {
		desc           string
		input          string
		dataToUpdate   []byte
		expectedOutput []byte
	}{
		{"Valid case", "123e4567-e89b-12d3-a456-426614174000",
			[]byte(fmt.Sprintf(`{"ID":"%v","Name":"jason","Dob":"12-06-2002","City":"Bangalore","Majors":"CSE","DId":1}`,
				"123e4567-e89b-12d3-a456-426614174000")),
			[]byte(fmt.Sprintf(`{"ID":"%v","Name":"jason","Dob":"12-06-2002","City":"Bangalore","Majors":"CSE","DId":1}`,
				"123e4567-e89b-12d3-a456-426614174000"))},
		{"for unmarshal error", "123e4567-e89b-12d3-a456-426614174000", []byte(nil), []byte("Invalid body")},
		{"id parsing error", "123e4567-e89b-12d3-a456426614174000", []byte(nil), []byte("Invalid id")},
		{"Invalid id", "123e4567-e89b-12d3-a456-426614171000", []byte(fmt.Sprintf(
			`{"ID":"%v","Name":"jason","Dob":"12-06-2002","City":"Bangalore","Majors":"CSE","DId":1}`,
			"123e4567-e89b-12d3-a456-426614174000")), []byte("Id not found")},
		{"InValid case", "123e4567-e89b-12d3-a456-426614174030",
			[]byte(fmt.Sprintf(`{"ID":"%v","Name":"jason","Dob":"12-06-2002","City":"Bangalore","Majors":"CSE","DId":1}`,
				"123e4567-e89b-12d3-a456-426614174000")),
			[]byte(`Id not found`)},
		{"InValid case", "123e4567-e89b-12d3-a456-426614174010",
			[]byte(fmt.Sprintf(`{"ID":"%v","Name":"jason","Dob":"12-06-2002","City":"Bangalore","Majors":"CSE","DId":1}`,
				"123e4567-e89b-12d3-a456-426614174000")),
			[]byte(`error`)},
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

func TestDeleteHandler(t *testing.T) {
	testcases := []struct {
		desc           string
		input          string
		expectedOutput int
	}{
		{"Valid ID", "123e4567-e89b-12d3-a456-426614174000", 204},
		{"ID not found", "123e4567-e89b-12d3-a456-426614174020", 404},
		{"invalid id", "123e4567-e89b-12d3-a456426614174020", 400},
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

func TestGetAll(t *testing.T) {
	testcases := []struct {
		desc              string
		name              string
		includeDepartment string
		expectedOutput    []byte
	}{
		{"valid name and including dept ", "jason", "true", []byte(fmt.Sprintf(
			`[{"ID":"%v","Name":"jason","Dob":"12-06-1998","City":"Bangalore","Majors":"CSE","Dept":{"ID":2,"Name":"TECH","FloorNo":2}}]`,
			"123e4567-e89b-12d3-a456-426614174000"))},
		{"valid name and not including dept ", "jason", "false", []byte(fmt.Sprintf(
			`[{"ID":"%v","Name":"jason","Dob":"12-06-1998","City":"Bangalore","Majors":"CSE","Dept":{"ID":0,"Name":"","FloorNo":0}}]`,
			"123e4567-e89b-12d3-a456-426614174000"))},
		{"invalid name and includeDepartment is true", "", "true",
			[]byte("error")},
		{"parsebool error", "jason", "abcd",
			[]byte("Invalid value for includeDept")},
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

func (m mockService) Create(ctx context.Context, employee *entities.Employee) (*entities.Employee, error) {
	name := employee.Name
	if name == "" {
		return &entities.Employee{}, errors.New("error")
	}

	id, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
	if err != nil {
		return &entities.Employee{}, errors.New("error")
	}

	return &entities.Employee{ID: id, Name: "jason", Dob: "12-06-2002", City: "Bangalore", Majors: "CSE", DId: 1}, nil
}

func (m mockService) Read(ctx context.Context, id uuid.UUID) (entities.EmployeeAndDepartment, error) {
	uid, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")

	if err != nil {
		return entities.EmployeeAndDepartment{}, err
	}

	if uid == id {
		return entities.EmployeeAndDepartment{ID: uid, Name: "jason", Dob: "12-06-2002", City: "Bangalore", Majors: "CSE",
			Dept: entities.Department{ID: 2, Name: "CSE", FloorNo: 2}}, nil
	} else if id == uuid.MustParse("123e4567-e89b-12d3-a456-426614174010") {
		return entities.EmployeeAndDepartment{ID: uid, Name: "jason", Dob: "12-06-2002", City: "Bangalore", Majors: "CSE",
			Dept: entities.Department{ID: 2, Name: "CSE", FloorNo: 2}}, nil
	}

	return entities.EmployeeAndDepartment{}, errors.New("error")
}

func (m mockService) Update(ctx context.Context, id uuid.UUID, employee *entities.Employee) (*entities.Employee, error) {
	uid, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
	if err != nil {
		return &entities.Employee{}, err
	}

	if uid == id {
		return &entities.Employee{ID: uid, Name: "jason", Dob: "12-06-2002", City: "Bangalore", Majors: "CSE", DId: 1}, nil
	}

	return &entities.Employee{}, errors.New("error")
}

func (m mockService) Delete(ctx context.Context, id uuid.UUID) (int, error) {
	uid, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
	if err != nil {
		return http.StatusBadRequest, nil
	}

	if id == uid {
		return http.StatusNoContent, nil
	}

	return http.StatusNotFound, &errorsHandler.IDNotFound{Msg: "not found"}
}

func (m mockService) ReadAll(para store.Parameters) ([]entities.EmployeeAndDepartment, error) {
	uid, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
	if err != nil {
		return []entities.EmployeeAndDepartment{}, err
	}

	if (para.Name != "") && (para.IncludeDepartment == true) {
		return []entities.EmployeeAndDepartment{{ID: uid, Name: "jason", Dob: "12-06-1998", City: "Bangalore",
			Majors: "CSE", Dept: entities.Department{ID: 2, Name: "TECH", FloorNo: 2}}}, nil
	} else if (para.Name != "") && (para.IncludeDepartment == false) {
		return []entities.EmployeeAndDepartment{{ID: uid, Name: "jason", Dob: "12-06-1998", City: "Bangalore",
			Majors: "CSE", Dept: entities.Department{}}}, nil
	}

	return []entities.EmployeeAndDepartment{}, errors.New("error")
}
