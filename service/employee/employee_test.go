package employee

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/google/uuid"

	entities "EmployeeDepartment/entities"
	"EmployeeDepartment/store"
)

type mockDatastore struct{}

func TestValidatePost(t *testing.T) {
	testcases := []struct {
		desc           string
		input          *entities.Employee
		expectedOutput *entities.Employee
	}{
		{desc: "invalid dob", input: &entities.Employee{uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			"jason", "12-06-1999", "Bangalore", "CSE", 2}, expectedOutput: &entities.Employee{}},
		{desc: "Valid input1", input: &entities.Employee{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			Name: "jason", Dob: "12-06-1998", City: "Bangalore", Majors: "MBA", DId: 1},
			expectedOutput: &entities.Employee{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), Name: "jason",
				Dob: "12-06-1998", City: "Bangalore", Majors: "MBA", DId: 1}},
		{desc: "Valid input2", input: &entities.Employee{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			Name: "jason", Dob: "12-06-1998", City: "Bangalore", Majors: "MCA", DId: 2},
			expectedOutput: &entities.Employee{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), Name: "jason",
				Dob: "12-06-1998", City: "Bangalore", Majors: "MCA", DId: 2}},
		{desc: "Valid input3", input: &entities.Employee{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			Name: "jason", Dob: "12-06-1998", City: "Bangalore", Majors: "CA", DId: 3},
			expectedOutput: &entities.Employee{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), Name: "jason",
				Dob: "12-06-1998", City: "Bangalore", Majors: "CA", DId: 3}},
		{desc: "Invalid deptid", input: &entities.Employee{uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			"jason", "12-06-1967", "Bangalore", "MCA", 1}, expectedOutput: &entities.Employee{}},
		{desc: "Invalid deptid", input: &entities.Employee{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			Name: "jason", Dob: "12-06-1967", City: "Bangalore", Majors: "CA", DId: 1}, expectedOutput: &entities.Employee{}},
		{desc: "Invalid deptid", input: &entities.Employee{uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			"jason", "12-06-1967", "Bangalore", "MBA", 2}, expectedOutput: &entities.Employee{}},
		{desc: "Invalid name", input: &entities.Employee{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			Dob: "12-06-1967", City: "Bangalore", Majors: "MBA", DId: 1}, expectedOutput: &entities.Employee{}},
		{desc: "Invalid deptid1", input: &entities.Employee{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			Name: "jason", Dob: "12-06-1967", City: "Bangalore", Majors: "MBA", DId: 5}, expectedOutput: &entities.Employee{}},
		{desc: "Invalid deptid2", input: &entities.Employee{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			Name: "jason", Dob: "12-06-1967", City: "Bangalore", Majors: "MB", DId: 5}, expectedOutput: &entities.Employee{}},
		{desc: "Invalid city", input: &entities.Employee{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			Name: "jason", Dob: "12-06-1968", City: "Delhi", Majors: "MCA", DId: 1}, expectedOutput: &entities.Employee{}},
		{desc: "Invalid majors", input: &entities.Employee{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			Name: "jason", Dob: "12-06-1968", City: "Kochi", Majors: "B.Sc", DId: 2}, expectedOutput: &entities.Employee{}},
		{desc: "Invalid date", input: &entities.Employee{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			Name: "jason", Dob: "13/06/1968", City: "Kochi", Majors: "B.Sc", DId: 2}, expectedOutput: &entities.Employee{}},
	}

	for i, tc := range testcases {
		a := New(mockDatastore{})

		actualOutput, _ := a.Create(context.TODO(), tc.input)
		if !reflect.DeepEqual(actualOutput, tc.expectedOutput) {
			t.Errorf("testcase %d %v failed got %v \n expected %v", i+1, tc.desc, actualOutput, tc.expectedOutput)
		}
	}
}

func TestValidatePut(t *testing.T) {
	testcases := []struct {
		desc           string
		id             uuid.UUID
		input          *entities.Employee
		expectedOutput *entities.Employee
	}{
		{
			desc: "Invalid dob",
			id:   uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			input: &entities.Employee{uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				"jason", "12-06-1999", "Bangalore", "CSE", 1},
			expectedOutput: &entities.Employee{},
		},
		{desc: "Valid input", id: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			input: &entities.Employee{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Name: "jason", Dob: "12-06-1998", City: "Bangalore", Majors: "CSE", DId: 2},
			expectedOutput: &entities.Employee{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Name: "jason", Dob: "12-06-1998", City: "Bangalore", Majors: "CSE", DId: 2}},
		{desc: "Invalid age", id: uuid.MustParse("123e4567-e89b-12d3-a456-426614174010"),
			input: &entities.Employee{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174010"), Name: "jason",
				Dob: "12-06-1999", City: "Bangalore", Majors: "CSE", DId: 2}, expectedOutput: &entities.Employee{}},
		{desc: "Invalid id", id: uuid.MustParse("123e4567-e89b-12d3-a456-426614104000"),
			input: &entities.Employee{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Name: "jason", Dob: "12-06-1974", City: "Bangalore", Majors: "CSE", DId: 2}, expectedOutput: &entities.Employee{}},
		{desc: "Invalid city", id: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			input: &entities.Employee{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Name: "jason", Dob: "12-06-1996", City: "Delhi", Majors: "MCA", DId: 2}, expectedOutput: &entities.Employee{}},
		{desc: "Invalid Majors", id: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			input: &entities.Employee{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Name: "jason", Dob: "12-06-1996", City: "Mysore", Majors: "B.Sc", DId: 2}, expectedOutput: &entities.Employee{}},
		{desc: "inValid department id", id: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			input: &entities.Employee{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Name: "jason", Dob: "12-06-1998", City: "Bangalore", Majors: "CSE", DId: 4},
			expectedOutput: nil},
	}
	for i, tc := range testcases {
		a := New(mockDatastore{})

		actualOutput, _ := a.Update(context.TODO(), tc.id, tc.input)

		if !reflect.DeepEqual(actualOutput, tc.expectedOutput) {
			t.Errorf("testcase %d failed got %v \n expected %v", i+1, actualOutput, tc.expectedOutput)
		}
	}
}

func TestValidateDelete(t *testing.T) {
	testcases := []struct {
		desc           string
		input          uuid.UUID
		expectedOutput int
	}{
		{"Valid ID", uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), 204},
		{"ID not found", uuid.MustParse("123e4567-e89b-12d3-a456-426614074000"), 404},
	}
	for i, tc := range testcases {
		a := New(mockDatastore{})

		actualOutput, _ := a.Delete(context.TODO(), tc.input)
		if actualOutput != tc.expectedOutput {
			t.Errorf("testcase %d failed got %v \n expected %v", i+1, actualOutput, tc.expectedOutput)
		}
	}
}

func TestGetById(t *testing.T) {
	testcases := []struct {
		desc           string
		id             uuid.UUID
		expectedOutput entities.EmployeeAndDepartment
	}{
		{desc: "valid name ", id: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			expectedOutput: entities.EmployeeAndDepartment{ID: uuid.MustParse(
				"123e4567-e89b-12d3-a456-426614174000"),
				Name: "jason", Dob: "12-06-1998", City: "Bangalore", Majors: "CSE",
				Dept: entities.Department{ID: 1, Name: "HR", FloorNo: 1}}},
		{"valid name ", uuid.MustParse("123e4567-e89b-12d3-a456-426614174020"),
			entities.EmployeeAndDepartment{}},
	}

	for i, tc := range testcases {
		a := New(mockDatastore{})

		actualOutput, _ := a.Read(context.TODO(), tc.id)
		if !reflect.DeepEqual(actualOutput, tc.expectedOutput) {
			t.Errorf("testcase %d failed got %v \n expected %v", i+1, actualOutput, tc.expectedOutput)
		}
	}
}

func TestGetAll(t *testing.T) {
	testcases := []struct {
		desc              string
		name              string
		includeDepartment bool
		expectedOutput    []entities.EmployeeAndDepartment
	}{
		{desc: "valid name and including dept ", name: "jason", includeDepartment: true,
			expectedOutput: []entities.EmployeeAndDepartment{{ID: uuid.MustParse(
				"123e4567-e89b-12d3-a456-426614174000"), Name: "jason", Dob: "12-06-1998", City: "Bangalore",
				Majors: "CSE", Dept: entities.Department{ID: 1, Name: "HR", FloorNo: 1}}}},
		{"valid name and not including dept ", "jason", false,
			[]entities.EmployeeAndDepartment{{ID: uuid.MustParse(
				"123e4567-e89b-12d3-a456-426614174000"), Name: "jason", Dob: "12-06-1998", City: "Bangalore",
				Majors: "CSE", Dept: entities.Department{}}}},
		{"invalid name and includeDepartment is true", "", true,
			[]entities.EmployeeAndDepartment{{}}},
	}

	for i, tc := range testcases {
		a := New(mockDatastore{})

		actualOutput, _ := a.ReadAll(store.Parameters{context.TODO(), tc.name,
			tc.includeDepartment})
		if !reflect.DeepEqual(actualOutput, tc.expectedOutput) {
			t.Errorf("testcase %d failed got %v \n expected %v", i+1, actualOutput, tc.expectedOutput)
		}
	}
}

func (m mockDatastore) Create(ctx context.Context, e *entities.Employee) (*entities.Employee, error) {
	if e.Name == "" {
		return &entities.Employee{}, errors.New("error")
	} else if e.DId == 2 {
		return &entities.Employee{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), Name: "jason",
			Dob: "12-06-1998", City: "Bangalore", Majors: "MCA", DId: 2}, nil
	} else if e.DId == 3 {
		return &entities.Employee{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), Name: "jason",
			Dob: "12-06-1998", City: "Bangalore", Majors: "CA", DId: 3}, nil
	}

	return &entities.Employee{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), Name: "jason",
		Dob: "12-06-1998", City: "Bangalore", Majors: "MBA", DId: 1}, nil
}

func (m mockDatastore) Update(ctx context.Context, id uuid.UUID, e *entities.Employee) (*entities.Employee, error) {
	if uuid.MustParse("123e4567-e89b-12d3-a456-426614174000") == id {
		return &entities.Employee{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), Name: "jason",
			Dob: "12-06-1998", City: "Bangalore", Majors: "CSE", DId: 2}, nil
	}

	return &entities.Employee{}, errors.New("error")
}

func (m mockDatastore) ReadDepartment(ctx context.Context, id int) (department entities.Department, err error) {
	if id == 1 {
		return entities.Department{ID: 1, Name: "HR", FloorNo: 1}, nil
	} else if id == 2 {
		return entities.Department{ID: 2, Name: "TECH", FloorNo: 2}, nil
	} else if id == 3 {
		return entities.Department{ID: 3, Name: "ACCOUNTS", FloorNo: 3}, nil
	}
	return entities.Department{}, errors.New("error")
}

func (m mockDatastore) Delete(ctx context.Context, id uuid.UUID) (int, error) {
	if id == uuid.MustParse("123e4567-e89b-12d3-a456-426614174000") {
		return http.StatusNoContent, nil
	}

	return http.StatusNotFound, errors.New("error")
}

func (m mockDatastore) Read(ctx context.Context, id uuid.UUID) (entities.EmployeeAndDepartment, error) {
	if id == uuid.MustParse("123e4567-e89b-12d3-a456-426614174000") {
		return entities.EmployeeAndDepartment{ID: id, Name: "jason", Dob: "12-06-1998", City: "Bangalore",
			Majors: "CSE", Dept: entities.Department{ID: 1, Name: "HR", FloorNo: 1}}, nil
	}

	return entities.EmployeeAndDepartment{}, errors.New("error")
}

func (m mockDatastore) ReadAll(para store.Parameters) ([]entities.EmployeeAndDepartment, error) {
	if (para.Name != "") && para.IncludeDepartment == true {
		out := make([]entities.EmployeeAndDepartment, 1)
		out[0] = entities.EmployeeAndDepartment{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			Name: "jason", Dob: "12-06-1998", City: "Bangalore", Majors: "CSE",
			Dept: entities.Department{ID: 1, Name: "HR", FloorNo: 1}}

		return out, nil
	} else if (para.Name != "") && (!para.IncludeDepartment) {
		out := make([]entities.EmployeeAndDepartment, 1)
		out[0] = entities.EmployeeAndDepartment{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			Name: "jason", Dob: "12-06-1998", City: "Bangalore", Majors: "CSE"}
		return out, nil
	}

	return []entities.EmployeeAndDepartment{{}}, errors.New("error")
}
