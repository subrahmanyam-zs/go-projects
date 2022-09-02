package department

import (
	"testing"

	"github.com/golang/mock/gomock"

	"developer.zopsmart.com/go/gofr/Emp-Dept/datastore"
	"developer.zopsmart.com/go/gofr/Emp-Dept/entities"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func TestPost(t *testing.T) {
	testCases := []struct {
		desc           string
		input          entities.Department
		expectedOutput interface{}
		err            error
	}{
		{desc: "Valid input", input: entities.Department{DeptID: 1, DeptName: "HR", FloorNo: 1},
			expectedOutput: entities.Department{DeptID: 1, DeptName: "HR", FloorNo: 1}, err: nil},
		{"already exists", entities.Department{1, "HR", 1}, nil, errors.EntityAlreadyExists{}},
	}

	for i, tc := range testCases {
		ctx := gofr.NewContext(nil, nil, gofr.New())
		ctrl := gomock.NewController(t)
		mockIndex := datastore.NewMockDepartment(ctrl)

		mockIndex.EXPECT().Post(ctx, tc.input).Return(tc.expectedOutput, tc.err)

		mockService := New(mockIndex)

		actualOutput, err := mockService.Post(ctx, tc.input)
		if actualOutput != tc.expectedOutput {
			t.Errorf("testcase %d failed got %v error %v\n expected %v error %v", i+1, actualOutput, err, tc.expectedOutput, tc.err)
		}
	}
}

func TestPostWithEdgeCases(t *testing.T) {
	testCases := []struct {
		desc           string
		input          entities.Department
		expectedOutput interface{}
		err            error
	}{
		{"Invalid floorNo", entities.Department{2, "TECH", 4}, nil, errors.InvalidParam{Param: []string{"floorNo"}}},
	}
	for i, tc := range testCases {
		ctx := gofr.NewContext(nil, nil, gofr.New())
		ctrl := gomock.NewController(t)
		mockIndex := datastore.NewMockDepartment(ctrl)

		mockService := Service{store: mockIndex}

		actualOutput, err := mockService.Post(ctx, tc.input)
		if actualOutput != tc.expectedOutput {
			t.Errorf("testcase %d failed got %v error %v\n expected %v error %v", i+1, actualOutput, err, tc.expectedOutput, tc.err)
		}
	}
}

func TestPut(t *testing.T) {
	testcases := []struct {
		desc           string
		id             int
		dataToUpdate   entities.Department
		expectedOutput interface{}
		err            error
	}{
		{
			"valid input",
			1,
			entities.Department{1, "HR", 1},
			entities.Department{1, "HR", 1}, nil,
		},
		{desc: "already exists", dataToUpdate: entities.Department{DeptID: 1, DeptName: "HR", FloorNo: 1}, err: errors.EntityAlreadyExists{}},
	}
	for i, tc := range testcases {
		ctx := gofr.NewContext(nil, nil, gofr.New())
		ctrl := gomock.NewController(t)
		mockIndex := datastore.NewMockDepartment(ctrl)

		mockIndex.EXPECT().Put(ctx, tc.id, tc.dataToUpdate).Return(tc.expectedOutput, tc.err)

		mockStore := Service{store: mockIndex}

		actualOutput, err := mockStore.Put(ctx, tc.id, tc.dataToUpdate)
		if actualOutput != tc.expectedOutput {
			t.Errorf("testcase %d failed got %v error %v \n expected %v error %v", i+1, actualOutput, err, tc.expectedOutput, tc.err)
		}
	}
}

func TestPutWithEdgeCases(t *testing.T) {
	testcases := []struct {
		desc           string
		id             int
		dataToUpdate   entities.Department
		expectedOutput interface{}
		err            error
	}{
		{desc: "Invalid floorNo", id: 2, dataToUpdate: entities.Department{DeptID: 2, DeptName: "THE", FloorNo: 3},
			err: errors.InvalidParam{Param: []string{"name"}}},
	}
	for i, tc := range testcases {
		ctx := gofr.NewContext(nil, nil, gofr.New())
		ctrl := gomock.NewController(t)
		mockIndex := datastore.NewMockDepartment(ctrl)

		mockStore := Service{store: mockIndex}

		actualOutput, err := mockStore.Put(ctx, tc.id, tc.dataToUpdate)
		if actualOutput != tc.expectedOutput {
			t.Errorf("testcase %d failed got %v error %v \n expected %v error %v", i+1, actualOutput, err, tc.expectedOutput, tc.err)
		}
	}
}

func TestDelete(t *testing.T) {
	testcases := []struct {
		desc           string
		id             int
		expectedOutput int
		err            error
	}{
		{"valid id", 1, 204, nil},
		{"id not found", 0, 404, errors.MissingParam{Param: []string{"id not found"}}},
	}
	for i, tc := range testcases {
		ctx := gofr.NewContext(nil, nil, gofr.New())
		ctrl := gomock.NewController(t)
		mockIndex := datastore.NewMockDepartment(ctrl)

		mockIndex.EXPECT().Delete(ctx, tc.id).Return(tc.expectedOutput, tc.err)

		mockStore := Service{store: mockIndex}

		actualOutput, err := mockStore.Delete(ctx, tc.id)
		if actualOutput != tc.expectedOutput {
			t.Errorf("testcase %d failed got %v error %v\n expected %v error %v", i+1, actualOutput, err, tc.expectedOutput, tc.err)
		}
	}
}
