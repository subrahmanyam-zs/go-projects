package department

import (
	"EmployeeDepartment/errorsHandler"
	"context"
	"errors"
	"net/http"
	"testing"

	"EmployeeDepartment/entities"
)

type mockDataStore struct{}

func TestValidatePost(t *testing.T) {
	testCases := []struct {
		desc           string
		input          entities.Department
		expectedOutput entities.Department
	}{
		{desc: "Valid input", input: entities.Department{ID: 1, Name: "HR", FloorNo: 1},
			expectedOutput: entities.Department{ID: 1, Name: "HR", FloorNo: 1}},
		{desc: "Invalid name", input: entities.Department{ID: 2, FloorNo: 2}},
		{desc: "Invalid floorNo", input: entities.Department{ID: 2, Name: "TECH"}},
		{desc: "Invalid id", input: entities.Department{Name: "TECH", FloorNo: 2}},
	}

	for i, tc := range testCases {
		a := New(mockDataStore{})

		actualOutput, _ := a.Create(context.TODO(), tc.input)
		if actualOutput != tc.expectedOutput {
			t.Errorf("testcase %d failed got %v \n expected %v", i+1, actualOutput, tc.expectedOutput)
		}
	}
}

func TestValidatePut(t *testing.T) {
	testcases := []struct {
		desc           string
		id             int
		dataToUpdate   entities.Department
		expectedOutput entities.Department
	}{
		{desc: "valid input", id: 1, dataToUpdate: entities.Department{ID: 1, Name: "HR", FloorNo: 1},
			expectedOutput: entities.Department{ID: 1, Name: "HR", FloorNo: 1}},
		{desc: "invalid id", dataToUpdate: entities.Department{ID: 1, Name: "HR", FloorNo: 1}},
		{desc: "Invalid name", id: 2, dataToUpdate: entities.Department{ID: 2, FloorNo: 2}},
		{desc: "Invalid floorNo", id: 2, dataToUpdate: entities.Department{ID: 2, Name: "TECH"}},
	}
	for i, tc := range testcases {
		a := New(mockDataStore{})

		actualOutput, _ := a.Update(context.TODO(), tc.id, tc.dataToUpdate)
		if actualOutput != tc.expectedOutput {
			t.Errorf("testcase %d failed got %v \n expected %v", i+1, actualOutput, tc.expectedOutput)
		}
	}
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

		actualOutput, _ := a.Delete(context.TODO(), tc.id)

		if actualOutput != tc.expectedOutput {
			t.Errorf("testcase %d failed got %v \n expected %v", i+1, actualOutput, tc.expectedOutput)
		}
	}
}

func TestGetDepartment(t *testing.T) {
	testcases := []struct {
		desc           string
		id             int
		expectedOutput entities.Department
	}{
		{"valid id", 1, entities.Department{1, "HR", 1}},
		{"Invalid id", 4, entities.Department{}},
	}
	for i, tc := range testcases {
		a := New(mockDataStore{})

		actualOutput, _ := a.GetDepartment(context.TODO(), tc.id)
		if actualOutput != tc.expectedOutput {
			t.Errorf("testcase %d failed got %v \n expected %v", i+1, actualOutput, tc.expectedOutput)
		}
	}
}

func (m mockDataStore) Create(ctx context.Context, department entities.Department) (entities.Department, error) {
	if (department.ID >= 1) && (department.ID <= 3) {
		return entities.Department{ID: 1, Name: "HR", FloorNo: 1}, nil
	}

	return entities.Department{}, errors.New("error")
}

func (m mockDataStore) Update(ctx context.Context, id int, department entities.Department) (entities.Department, error) {
	if (id >= 1) && (id <= 3) {
		return department, nil
	}

	return entities.Department{}, errors.New("error")
}

func (m mockDataStore) Delete(ctx context.Context, id int) (int, error) {
	if (id >= 1) && (id <= 3) {
		return http.StatusNoContent, nil
	}

	return http.StatusNotFound, errors.New("error")
}
func (m mockDataStore) GetDepartment(ctx context.Context, id int) (entities.Department, error) {
	if id == 1 {
		return entities.Department{1, "HR", 1}, nil
	}
	return entities.Department{}, errorsHandler.IDNotFound{Msg: "Id not found"}
}
