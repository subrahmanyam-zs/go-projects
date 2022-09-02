package employee

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	uuidMongo "go.mongodb.org/mongo-driver/x/mongo/driver/uuid"

	"developer.zopsmart.com/go/gofr/Emp-Dept/datastore"
	"developer.zopsmart.com/go/gofr/Emp-Dept/entities"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func TestPost(t *testing.T) {
	testcases := []struct {
		desc           string
		input          entities.Employee
		expectedOutput entities.Employee
		dept           entities.Department
		err            error
	}{
		{
			"valid input",
			entities.Employee{uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				"jason", "12-06-1998", "Bangalore", "CSE", 2},
			entities.Employee{uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				"jason", "12-06-1998", "Bangalore", "CSE", 2}, entities.Department{2, "TECH", 2},
			nil,
		},
		{
			"valid input",
			entities.Employee{uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				"jason", "12-06-1998", "Bangalore", "CA", 3},
			entities.Employee{uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				"jason", "12-06-1998", "Bangalore", "CA", 3}, entities.Department{3, "ACCOUNTS", 3},
			nil,
		},
		{
			"valid input",
			entities.Employee{uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				"jason", "12-06-1998", "Bangalore", "MBA", 1},
			entities.Employee{uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				"jason", "12-06-1998", "Bangalore", "MBA", 1}, entities.Department{1, "HR", 1},
			nil,
		},
		{
			"Already Exists",
			entities.Employee{uuid.UUID(uuidMongo.UUID(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))),
				"jason", "12-06-1998", "Bangalore", "CSE", 2},
			entities.Employee{}, entities.Department{2, "TECH", 2}, errors.EntityAlreadyExists{},
		},
		{desc: "Valid input2", input: entities.Employee{ID: uuid.UUID(uuidMongo.UUID(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))),
			Name: "jason", Dob: "12-06-1998", City: "Bangalore", Majors: "MCA", DeptID: 2},
			expectedOutput: entities.Employee{ID: uuid.UUID(uuidMongo.UUID(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))), Name: "jason",
				Dob: "12-06-1998", City: "Bangalore", Majors: "MCA", DeptID: 2}, dept: entities.Department{2, "TECH", 2}, err: nil},
	}

	for i, tc := range testcases {
		ctx := gofr.NewContext(nil, nil, gofr.New())
		ctrl := gomock.NewController(t)
		mockIndex := datastore.NewMockEmployee(ctrl)

		mockIndex.EXPECT().GetDepartment(ctx, tc.input.DeptID).Return(tc.dept, nil)
		mockIndex.EXPECT().Post(ctx, tc.input).Return(tc.expectedOutput, tc.err)

		mockStore := Service{dataStore: mockIndex}
		actualOutput, err := mockStore.Post(ctx, tc.input)
		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("testcase %d %v failed got %v error %v \n expected %v error %v", i+1, tc.desc, actualOutput, err, tc.expectedOutput, tc.err)
		}
	}
}

func TestPostWithEdgeCases(t *testing.T) {
	testcases := []struct {
		desc           string
		input          entities.Employee
		expectedOutput interface{}
		dept           entities.Department
		err            error
	}{
		{
			"Invalid majors",
			entities.Employee{uuid.UUID(uuidMongo.UUID(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))),
				"jason", "12-06-1998", "Bangalore", "CE", 2},
			nil, entities.Department{},
			nil,
		},
		{
			"Invalid majors",
			entities.Employee{uuid.UUID(uuidMongo.UUID(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))),
				"jason", "12-06-1998", "Bangalore", "CE", 2},
			nil, entities.Department{},
			errors.InvalidParam{Param: []string{"invalid"}},
		},
		{
			"Invalid deptid",
			entities.Employee{uuid.UUID(uuidMongo.UUID(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))),
				"jason", "12-06-1967", "Bangalore", "MCA", 1},
			nil, entities.Department{1, "HR", 1},
			nil,
		},
		{
			"Invalid deptid",
			entities.Employee{uuid.UUID(uuidMongo.UUID(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))),
				"jason", "12-06-1967", "Bangalore", "CA", 1},
			nil, entities.Department{1, "HR", 1},
			nil,
		},

		{
			"Invalid deptid",
			entities.Employee{uuid.UUID(uuidMongo.UUID(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))),
				"jason", "12-06-1967", "Bangalore", "MBA", 2},
			nil, entities.Department{2, "TECH", 2}, nil,
		},
		{desc: "Invalid name", input: entities.Employee{
			uuid.UUID(uuidMongo.UUID(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))), "", "12-06-1967", "Bangalore", "MBA", 0,
		}, expectedOutput: nil, dept: entities.Department{}, err: nil},
		{desc: "Invalid deptid1", input: entities.Employee{
			uuid.UUID(uuidMongo.UUID(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))), "jason", "12-06-1967", "Bangalore", "MBA", 0,
		}, expectedOutput: nil, dept: entities.Department{}, err: nil},
		{desc: "Invalid deptid2", input: entities.Employee{ID: uuid.UUID(uuidMongo.UUID(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))),
			Name: "jason", Dob: "12-06-1967", City: "Bangalore", Majors: "MB", DeptID: 5}, dept: entities.Department{},
			expectedOutput: nil, err: nil},
		{desc: "Invalid city", input: entities.Employee{
			uuid.UUID(uuidMongo.UUID(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))),
			"jason", "12-06-1968", "Delhi", "MCA", 0,
		}, expectedOutput: nil, dept: entities.Department{}, err: nil},
		{desc: "Invalid majors", input: entities.Employee{
			uuid.UUID(uuidMongo.UUID(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))),
			"jason",
			"12-06-1968",
			"Kochi",
			"B.Sc",
			0,
		}, expectedOutput: nil, err: nil},
		{
			"Invalid date",
			entities.Employee{ID: uuid.UUID(uuidMongo.UUID(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))),
				Name: "jason", Dob: "13/061968", City: "Kochi", Majors: "B.Sc", DeptID: 2},
			nil, entities.Department{},
			nil,
		},
	}

	for i, tc := range testcases {
		ctx := gofr.NewContext(nil, nil, gofr.New())
		ctrl := gomock.NewController(t)
		mockIndex := datastore.NewMockEmployee(ctrl)

		mockIndex.EXPECT().GetDepartment(ctx, tc.input.DeptID).Return(tc.dept, tc.err)
		mockStore := New(mockIndex)

		actualOutput, err := mockStore.Post(ctx, tc.input)
		if !reflect.DeepEqual(actualOutput, tc.expectedOutput) {
			t.Errorf("testcase %d %v failed got %v error %v \n expected %v error %v", i+1, tc.desc, actualOutput, err, tc.expectedOutput, tc.err)
		}
	}
}

func TestPut(t *testing.T) {
	testcases := []struct {
		desc           string
		id             uuid.UUID
		dataToUpdate   entities.Employee
		expectedOutput interface{}
		err            error
	}{
		{
			"valid input", uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			entities.Employee{uuid.UUID(uuidMongo.UUID(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))),
				"jason", "12-06-1998", "Bangalore", "CSE", 2},
			entities.Employee{uuid.UUID(uuidMongo.UUID(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))),
				"jason", "12-06-1998", "Bangalore", "CSE", 2},
			nil,
		},
		{
			"Already Exists", uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			entities.Employee{uuid.UUID(uuidMongo.UUID(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))),
				"jason", "12-06-1998", "Bangalore", "CSE", 2},
			nil, errors.EntityAlreadyExists{},
		},
	}

	for i, tc := range testcases {
		ctx := gofr.NewContext(nil, nil, gofr.New())
		ctrl := gomock.NewController(t)
		mockIndex := datastore.NewMockEmployee(ctrl)

		mockIndex.EXPECT().GetDepartment(ctx, tc.dataToUpdate.DeptID).Return(entities.Department{2, "Tech", 2}, nil)
		mockIndex.EXPECT().Put(ctx, tc.id, tc.dataToUpdate).Return(tc.expectedOutput, tc.err)

		mockStore := Service{dataStore: mockIndex}

		actualOutput, err := mockStore.Put(ctx, tc.id, tc.dataToUpdate)
		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("testcase %d %v failed got %v error %v \n expected %v error %v", i+1, tc.desc, actualOutput, err, tc.expectedOutput, tc.err)
		}
	}
}

func TestPutWithEdgeCases(t *testing.T) {
	testcases := []struct {
		desc           string
		id             uuid.UUID
		dataToUpdate   entities.Employee
		expectedOutput interface{}
		err            error
	}{
		{
			"Invalid majors", uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			entities.Employee{uuid.UUID(uuidMongo.UUID(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))),
				"jason", "12-06-1998", "Bangaloe", "CSE", 2},
			nil,
			nil,
		},
		{
			"Already Exists", uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			entities.Employee{uuid.UUID(uuidMongo.UUID(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))),
				"jason", "12-06-1998", "Bangalore", "CSE", 2},
			nil, errors.EntityAlreadyExists{},
		},
	}

	for i, tc := range testcases {
		ctx := gofr.NewContext(nil, nil, gofr.New())
		ctrl := gomock.NewController(t)
		mockIndex := datastore.NewMockEmployee(ctrl)

		mockIndex.EXPECT().GetDepartment(ctx, tc.dataToUpdate.DeptID).Return(entities.Department{2, "Tech", 2}, tc.err)

		mockStore := Service{dataStore: mockIndex}

		actualOutput, err := mockStore.Put(ctx, tc.id, tc.dataToUpdate)
		if !reflect.DeepEqual(actualOutput, tc.expectedOutput) {
			t.Errorf("testcase %d %v failed got %v error %v \n expected %v error %v", i+1, tc.desc, actualOutput, err, tc.expectedOutput, tc.err)
		}
	}
}

func TestValidateDelete(t *testing.T) {
	testcases := []struct {
		desc           string
		input          uuid.UUID
		expectedOutput int
		err            error
	}{
		{"Valid ID", uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), 204, nil},
		{"ID not found", uuid.MustParse("123e4567-e89b-12d3-a456-426614074000"),
			404, errors.EntityNotFound{ID: "123e4567-e89b-12d3-a456-426614074000"}},
	}
	for i, tc := range testcases {
		ctx := gofr.NewContext(nil, nil, gofr.New())
		ctrl := gomock.NewController(t)
		mockIndex := datastore.NewMockEmployee(ctrl)

		mockIndex.EXPECT().Delete(ctx, tc.input).Return(tc.expectedOutput, tc.err)

		mockStore := Service{dataStore: mockIndex}

		actualOutput, err := mockStore.Delete(ctx, tc.input)
		if actualOutput != tc.expectedOutput {
			t.Errorf("testcase %d failed got %v error %v \n expected %v error %v", i+1, actualOutput, err, tc.expectedOutput, tc.err)
		}
	}
}

func TestGet(t *testing.T) {
	testcases := []struct {
		desc           string
		id             uuid.UUID
		expectedOutput interface{}
		err            error
	}{
		{
			"valid id ",
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			entities.EmpDept{uuid.UUID(uuidMongo.UUID(
				uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))), "jason", "12-06-1998", "Kochi", "CSE",
				entities.Department{DeptID: 1, DeptName: "HR", FloorNo: 1}},
			nil,
		},
		{"id not found ", uuid.MustParse("123e4567-e89b-12d3-a456-426614174020"),
			nil, errors.EntityNotFound{ID: "123e4567-e89b-12d3-a456-426614174000"}},
	}

	for i, tc := range testcases {
		ctx := gofr.NewContext(nil, nil, gofr.New())
		ctrl := gomock.NewController(t)
		mockIndex := datastore.NewMockEmployee(ctrl)

		mockIndex.EXPECT().Get(ctx, tc.id).Return(tc.expectedOutput, tc.err)

		mockStore := Service{dataStore: mockIndex}

		actualOutput, err := mockStore.Get(ctx, tc.id)
		if !reflect.DeepEqual(actualOutput, tc.expectedOutput) {
			t.Errorf("testcase %d failed got %v error %v\n expected %v error %v", i+1, actualOutput, err, tc.expectedOutput, tc.err)
		}
	}
}

func TestGetAll(t *testing.T) {
	testcases := []struct {
		desc              string
		name              string
		includeDepartment bool
		expectedOutput    interface{}
		err               error
	}{
		{"valid name and including dept ", "jason", true, []entities.EmpDept{{
			uuid.UUID(uuidMongo.UUID(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))), "jason", "12-06-1998", "Kochi", "CSE",
			entities.Department{
				2, "TECH", 2}}}, nil},
		{"name not found ", "roy", true, nil, errors.EntityNotFound{ID: "roy"}},
	}

	for i, tc := range testcases {
		ctx := gofr.NewContext(nil, nil, gofr.New())
		ctrl := gomock.NewController(t)
		mockIndex := datastore.NewMockEmployee(ctrl)

		mockIndex.EXPECT().GetAll(ctx, tc.name, tc.includeDepartment).Return(tc.expectedOutput, tc.err)

		mockStore := Service{dataStore: mockIndex}

		actualOutput, _ := mockStore.GetAll(ctx, tc.name, tc.includeDepartment)
		if !reflect.DeepEqual(actualOutput, tc.expectedOutput) {
			t.Errorf("testcase %d failed got %v \n expected %v", i+1, actualOutput, tc.expectedOutput)
		}
	}
}
