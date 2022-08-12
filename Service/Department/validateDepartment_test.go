package Department

import (
	"EmployeeDepartment/Handler/Entities"
	"errors"
	"net/http"
	"testing"
)

func TestValidatePost(t *testing.T) {
	testCases := []struct {
		desc           string
		input          Entities.Department
		expectedOutput Entities.Department
	}{
		{"Valid input", Entities.Department{1, "HR", 1}, Entities.Department{1, "HR", 1}},
		{"Invalid name", Entities.Department{2, "", 2}, Entities.Department{}},
		{"Invalid floorNo", Entities.Department{2, "TECH", 0}, Entities.Department{}},
		{"Invalid id", Entities.Department{0, "TECH", 2}, Entities.Department{}},
	}

	for i, tc := range testCases {
		a := New(mockDataStore{})
		actualOutput := a.validatePost(tc.input)
		if actualOutput != tc.expectedOutput {
			t.Errorf("testcase %d failed got %v \n expected %v", i+1, actualOutput, tc.expectedOutput)
		}
	}

}

func (m mockDataStore) Create(department Entities.Department) (Entities.Department, error) {
	if (department.Id >= 1) && (department.Id <= 3) {
		return Entities.Department{1, "HR", 1}, nil
	}
	return Entities.Department{}, errors.New("error")
}

func TestValidatePut(t *testing.T) {
	testcases := []struct {
		desc           string
		id             int
		dataToUpdate   Entities.Department
		expectedOutput Entities.Department
	}{
		{"valid input", 1, Entities.Department{1, "HR", 1}, Entities.Department{1, "HR", 1}},
		{"invalid id", 0, Entities.Department{1, "HR", 1}, Entities.Department{}},
		{"Invalid name", 2, Entities.Department{2, "", 2}, Entities.Department{}},
		{"Invalid floorNo", 2, Entities.Department{2, "TECH", 0}, Entities.Department{}},
	}
	for i, tc := range testcases {
		a := New(mockDataStore{})
		actualOutput := a.validatePut(tc.id, tc.dataToUpdate)
		if actualOutput != tc.expectedOutput {
			t.Errorf("testcase %d failed got %v \n expected %v", i+1, actualOutput, tc.expectedOutput)
		}
	}
}

func (m mockDataStore) Update(id int, department Entities.Department) (Entities.Department, error) {
	if (id >= 1) && (id <= 3) {
		return department, nil
	}
	return Entities.Department{}, errors.New("error")
}

func TestValidateDelete(t *testing.T) {
	testcases := []struct {
		desc           string
		id             int
		expectedOutput int
	}{
		{"valid id", 1, 204},
		{"invalid id", 0, 404},
	}
	for i, tc := range testcases {
		a := New(mockDataStore{})

		actualOutput := a.validateDelete(tc.id)

		if actualOutput != tc.expectedOutput {
			t.Errorf("testcase %d failed got %v \n expected %v", i+1, actualOutput, tc.expectedOutput)
		}
	}
}

func (m mockDataStore) Delete(id int) (int, error) {
	if (id >= 1) && (id <= 3) {
		return http.StatusNoContent, nil
	}
	return http.StatusNotFound, errors.New("error")
}

type mockDataStore struct{}
