package department

import (
	"EmployeeDepartment/entities"
	"errors"
	"net/http"
	"testing"
)

func TestValidatePost(t *testing.T) {
	testCases := []struct {
		desc           string
		input          entities.Department
		expectedOutput entities.Department
	}{
		{"Valid input", entities.Department{1, "HR", 1}, entities.Department{1, "HR", 1}},
		{"Invalid name", entities.Department{2, "", 2}, entities.Department{}},
		{"Invalid floorNo", entities.Department{2, "TECH", 0}, entities.Department{}},
		{"Invalid id", entities.Department{0, "TECH", 2}, entities.Department{}},
	}

	for i, tc := range testCases {
		a := New(mockDataStore{})
		actualOutput := a.validatePost(tc.input)
		if actualOutput != tc.expectedOutput {
			t.Errorf("testcase %d failed got %v \n expected %v", i+1, actualOutput, tc.expectedOutput)
		}
	}

}

func (m mockDataStore) Create(department entities.Department) (entities.Department, error) {
	if (department.Id >= 1) && (department.Id <= 3) {
		return entities.Department{1, "HR", 1}, nil
	}
	return entities.Department{}, errors.New("error")
}

func TestValidatePut(t *testing.T) {
	testcases := []struct {
		desc           string
		id             int
		dataToUpdate   entities.Department
		expectedOutput entities.Department
	}{
		{"valid input", 1, entities.Department{1, "HR", 1}, entities.Department{1, "HR", 1}},
		{"invalid id", 0, entities.Department{1, "HR", 1}, entities.Department{}},
		{"Invalid name", 2, entities.Department{2, "", 2}, entities.Department{}},
		{"Invalid floorNo", 2, entities.Department{2, "TECH", 0}, entities.Department{}},
	}
	for i, tc := range testcases {
		a := New(mockDataStore{})
		actualOutput := a.validatePut(tc.id, tc.dataToUpdate)
		if actualOutput != tc.expectedOutput {
			t.Errorf("testcase %d failed got %v \n expected %v", i+1, actualOutput, tc.expectedOutput)
		}
	}
}

func (m mockDataStore) Update(id int, department entities.Department) (entities.Department, error) {
	if (id >= 1) && (id <= 3) {
		return department, nil
	}
	return entities.Department{}, errors.New("error")
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
