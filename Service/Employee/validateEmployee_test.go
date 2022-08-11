package Employee

import (
	"EmployeeDepartment/Handler/Entities"
	"errors"
	"github.com/google/uuid"
	"net/http"
	"reflect"
	"testing"
)

var uid = uuid.New()

func TestValidatePost(t *testing.T) {
	testcases := []struct {
		desc           string
		input          Entities.Employee
		expectedOutput Entities.Employee
	}{
		{"Valid input", Entities.Employee{uid, "jason", "12-06-1999", "Bangalore", "CSE", 1}, Entities.Employee{}},
		{"Valid input", Entities.Employee{uid, "jason", "12-06-1998", "Bangalore", "CSE", 1},
			Entities.Employee{uid, "jason", "12-06-1998", "Bangalore", "CSE", 1}},
	}

	for i, tc := range testcases {
		a := New(mockDatastore{})
		actualOutput := a.validatePost(tc.input)
		if actualOutput != tc.expectedOutput {
			t.Errorf("testcase %d failed got %v \n expected %v", i+1, actualOutput, tc.expectedOutput)

		}
	}
}
func (m mockDatastore) Create(e Entities.Employee) (Entities.Employee, error) {
	if e.Name == "" {
		return Entities.Employee{}, errors.New("error")
	}
	return Entities.Employee{uid, "jason", "12-06-1998", "Bangalore", "CSE", 1}, nil
}

func TestValidatePut(t *testing.T) {
	testcases := []struct {
		desc           string
		id             string
		input          Entities.Employee
		expectedOutput Entities.Employee
	}{
		{"Invalid age", uid.String(), Entities.Employee{uid, "jason", "12-06-1999", "Bangalore", "CSE", 1}, Entities.Employee{}},
		{"Valid input", uid.String(), Entities.Employee{uid, "jason", "12-06-1998", "Bangalore", "CSE", 1},
			Entities.Employee{uid, "jason", "12-06-1998", "Bangalore", "CSE", 1}},
		{"Invalid age", "123e4567-e89b-12d3-a456-426614174000", Entities.Employee{uid, "jason", "12-06-1999", "Bangalore", "CSE", 1},
			Entities.Employee{}},
	}

	for i, tc := range testcases {
		a := New(mockDatastore{})

		actualOutput := a.validatePut(tc.id, tc.input)

		if actualOutput != tc.expectedOutput {
			t.Errorf("testcase %d failed got %v \n expected %v", i+1, actualOutput, tc.expectedOutput)
		}
	}
}

func (m mockDatastore) Update(id string, e Entities.Employee) (Entities.Employee, error) {
	if uid.String() == id {
		return Entities.Employee{uid, "jason", "12-06-1998", "Bangalore", "CSE", 1}, nil
	}
	return Entities.Employee{}, errors.New("error")
}

func TestValidateDelete(t *testing.T) {
	testcases := []struct {
		desc           string
		input          string
		expectedOutput int
	}{
		{"Valid Id", uid.String(), 204},
		{"Id not found", "123e4567-e89b-12d3-a456-426614174000", 404},
	}
	for i, tc := range testcases {
		a := New(mockDatastore{})

		actualOutput := a.validateDelete(tc.input)
		if actualOutput != tc.expectedOutput {
			t.Errorf("testcase %d failed got %v \n expected %v", i+1, actualOutput, tc.expectedOutput)
		}
	}
}

func (m mockDatastore) Delete(id uuid.UUID) (int, error) {
	if id == uid {
		return http.StatusNoContent, nil
	}
	return http.StatusNotFound, errors.New("error")
}

func TestGetById(t *testing.T) {
	testcases := []struct {
		desc           string
		input          string
		expectedOutput Entities.Employee
	}{
		{"Valid Input", uid.String(), Entities.Employee{uid, "jason", "12-06-1998", "Bangalore", "CSE", 1}},
		//{"Invalid Id", "00000-12223-122-1222323-2133", Entities.Employee{}},
		{"Id not found", "123e4567-e89b-12d3-a456-426614174000", Entities.Employee{}},
	}
	for i, tc := range testcases {
		a := New(mockDatastore{})
		actualOutput := a.validateGetById(tc.input)
		if actualOutput != tc.expectedOutput {
			t.Errorf("testcase %d failed got %v \n expected %v", i+1, actualOutput, tc.expectedOutput)
		}
	}
}

func (m mockDatastore) Read(id uuid.UUID) (Entities.Employee, error) {
	if id == uid {
		return Entities.Employee{uid, "jason", "12-06-1998", "Bangalore", "CSE", 1}, nil
	}
	return Entities.Employee{}, errors.New("error")
}

func TestGetAll(t *testing.T) {
	testcases := []struct {
		desc              string
		name              string
		includeDepartment bool
		expectedOutput    []Entities.EmployeeAndDepartment
	}{
		{"valid name and including dept ", "jason", true,
			[]Entities.EmployeeAndDepartment{{uid.String(), "jason", "12-06-1998", "Bangalore", "CSE", Entities.Department{1, "HR", 1}}}},
		{"valid name and not including dept ", "jason", false,
			[]Entities.EmployeeAndDepartment{{uid.String(), "jason", "12-06-1998", "Bangalore", "CSE", Entities.Department{}}}},
		{"invalid name and includeDepartment is true", "", true,
			[]Entities.EmployeeAndDepartment{{}}},
	}

	for i, tc := range testcases {
		a := New(mockDatastore{})

		actualOutput := a.validateGetAll(tc.name, tc.includeDepartment)
		if !reflect.DeepEqual(actualOutput, tc.expectedOutput) {
			t.Errorf("testcase %d failed got %v \n expected %v", i+1, actualOutput, tc.expectedOutput)

		}
	}
}

func (m mockDatastore) ReadAll(name string, includeDepartment bool) ([]Entities.EmployeeAndDepartment, error) {
	if (name != "") && includeDepartment == true {
		out := make([]Entities.EmployeeAndDepartment, 1, 1)
		out[0] = Entities.EmployeeAndDepartment{uid.String(), "jason", "12-06-1998", "Bangalore", "CSE", Entities.Department{1, "HR", 1}}
		return out, nil
	} else if (name != "") && (!includeDepartment) {
		out := make([]Entities.EmployeeAndDepartment, 1, 1)
		out[0] = Entities.EmployeeAndDepartment{uid.String(), "jason", "12-06-1998", "Bangalore", "CSE", Entities.Department{}}
		return out, nil
	}
	return []Entities.EmployeeAndDepartment{{}}, errors.New("error")
}

type mockDatastore struct{}
