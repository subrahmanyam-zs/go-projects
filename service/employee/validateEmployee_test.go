package employee

import (
	entities2 "EmployeeDepartment/entities"
	"errors"
	"github.com/google/uuid"
	"net/http"
	"reflect"
	"testing"
)

var uid = uuid.New()
var invalidUid = uuid.New()

func TestValidatePost(t *testing.T) {
	testcases := []struct {
		desc           string
		input          entities2.Employee
		expectedOutput entities2.Employee
	}{
		{"Valid input", entities2.Employee{uid, "jason", "12-06-1999", "Bangalore", "CSE", 2}, entities2.Employee{}},
		{"Valid input", entities2.Employee{uid, "jason", "12-06-1998", "Bangalore", "MBA", 1},
			entities2.Employee{uid, "jason", "12-06-1998", "Bangalore", "MBA", 1}},
		{"Valid input", entities2.Employee{uid, "jason", "12-06-1999", "Bangalore", "MCA", 2}, entities2.Employee{}},
	}

	for i, tc := range testcases {
		a := New(mockDatastore{})
		actualOutput, _ := a.Create(tc.input)
		if actualOutput != tc.expectedOutput {
			t.Errorf("testcase %d failed got %v \n expected %v", i+1, actualOutput, tc.expectedOutput)

		}
	}
}
func (m mockDatastore) Create(e entities2.Employee) (entities2.Employee, error) {
	if e.Name == "" {
		return entities2.Employee{}, errors.New("error")
	}
	return entities2.Employee{uid, "jason", "12-06-1998", "Bangalore", "MBA", 1}, nil
}

func TestValidatePut(t *testing.T) {
	testcases := []struct {
		desc           string
		id             uuid.UUID
		input          entities2.Employee
		expectedOutput entities2.Employee
	}{
		{"Invalid age", uid, entities2.Employee{uid, "jason", "12-06-1999", "Bangalore", "CSE", 1}, entities2.Employee{}},
		{"Valid input", uid, entities2.Employee{uid, "jason", "12-06-1998", "Bangalore", "CSE", 1},
			entities2.Employee{uid, "jason", "12-06-1998", "Bangalore", "CSE", 1}},
		{"Invalid age", invalidUid, entities2.Employee{uid, "jason", "12-06-1999", "Bangalore", "CSE", 1},
			entities2.Employee{}},
	}

	for i, tc := range testcases {
		a := New(mockDatastore{})

		actualOutput, _ := a.Update(tc.id, tc.input)

		if actualOutput != tc.expectedOutput {
			t.Errorf("testcase %d failed got %v \n expected %v", i+1, actualOutput, tc.expectedOutput)
		}
	}
}

func (m mockDatastore) Update(id uuid.UUID, e entities2.Employee) (entities2.Employee, error) {
	if uid == id {
		return entities2.Employee{uid, "jason", "12-06-1998", "Bangalore", "CSE", 1}, nil
	}
	return entities2.Employee{}, errors.New("error")
}

func TestValidateDelete(t *testing.T) {
	testcases := []struct {
		desc           string
		input          uuid.UUID
		expectedOutput int
	}{
		{"Valid Id", uid, 204},
		{"Id not found", invalidUid, 404},
	}
	for i, tc := range testcases {
		a := New(mockDatastore{})

		actualOutput, _ := a.Delete(tc.input)
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
		id             uuid.UUID
		expectedOutput entities2.EmployeeAndDepartment
	}{
		{"valid name ", uid,
			entities2.EmployeeAndDepartment{uid.String(), "jason", "12-06-1998", "Bangalore", "CSE", entities2.Department{1, "HR", 1}}},
		//{"invalid name and includeDepartment is true", "", []entities.EmployeeAndDepartment{{}}},
	}

	for i, tc := range testcases {
		a := New(mockDatastore{})

		actualOutput, _ := a.Read(tc.id)
		if !reflect.DeepEqual(actualOutput, tc.expectedOutput) {
			t.Errorf("testcase %d failed got %v \n expected %v", i+1, actualOutput, tc.expectedOutput)

		}
	}
}
func (m mockDatastore) Read(id uuid.UUID) (entities2.EmployeeAndDepartment, error) {
	if id == uid {
		//out := make([]entities2.EmployeeAndDepartment, 1, 1)
		return entities2.EmployeeAndDepartment{id.String(), "jason", "12-06-1997", "Bangalore", "CSE", entities2.Department{2, "TECH", 2}}, nil
	}
	return entities2.EmployeeAndDepartment{}, errors.New("error")
}

func TestGetAll(t *testing.T) {
	testcases := []struct {
		desc              string
		name              string
		includeDepartment bool
		expectedOutput    []entities2.EmployeeAndDepartment
	}{
		{"valid name and including dept ", "jason", true,
			[]entities2.EmployeeAndDepartment{{uid.String(), "jason", "12-06-1998", "Bangalore", "CSE", entities2.Department{1, "HR", 1}}}},
		{"valid name and not including dept ", "jason", false,
			[]entities2.EmployeeAndDepartment{{uid.String(), "jason", "12-06-1998", "Bangalore", "CSE", entities2.Department{}}}},
		{"invalid name and includeDepartment is true", "", true,
			[]entities2.EmployeeAndDepartment{{}}},
	}

	for i, tc := range testcases {
		a := New(mockDatastore{})

		actualOutput, _ := a.ReadAll(tc.name, tc.includeDepartment)
		if !reflect.DeepEqual(actualOutput, tc.expectedOutput) {
			t.Errorf("testcase %d failed got %v \n expected %v", i+1, actualOutput, tc.expectedOutput)

		}
	}
}

func (m mockDatastore) ReadAll(name string, includeDepartment bool) ([]entities2.EmployeeAndDepartment, error) {
	if (name != "") && includeDepartment == true {
		out := make([]entities2.EmployeeAndDepartment, 1, 1)
		out[0] = entities2.EmployeeAndDepartment{uid.String(), "jason", "12-06-1998", "Bangalore", "CSE", entities2.Department{1, "HR", 1}}
		return out, nil
	} else if (name != "") && (!includeDepartment) {
		out := make([]entities2.EmployeeAndDepartment, 1, 1)
		out[0] = entities2.EmployeeAndDepartment{uid.String(), "jason", "12-06-1998", "Bangalore", "CSE", entities2.Department{}}
		return out, nil
	}
	return []entities2.EmployeeAndDepartment{{}}, errors.New("error")
}

func (m mockDatastore) ReadDepartment(id int) (department entities2.Department, err error) {
	if id == 1 {
		return entities2.Department{1, "HR", 1}, nil
	} else if id == 2 {
		return entities2.Department{2, "TECH", 2}, nil
	} else if id == 3 {
		return entities2.Department{3, "ACCOUNTS", 3}, nil
	}
	return entities2.Department{}, errors.New("error")
}

type mockDatastore struct{}
