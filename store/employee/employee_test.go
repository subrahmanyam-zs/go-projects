package employee

import (
	"EmployeeDepartment/errorsHandler"
	"context"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"reflect"
	"testing"

	entities "EmployeeDepartment/entities"
	"EmployeeDepartment/store"
)

func TestCreateEmployee(t *testing.T) {
	testcases := []struct {
		desc           string
		input          *entities.Employee
		rowsAffec      int64
		expectedOutput *entities.Employee
		err            error
	}{
		{desc: "valid input", input: &entities.Employee{Name: "jason", Dob: "12-06-1998", City: "Bangalore",
			Majors: "MBA", DId: 1}, rowsAffec: 1,
			expectedOutput: &entities.Employee{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Name: "jason", Dob: "12-06-1998", City: "Bangalore", Majors: "MBA", DId: 1}},
		{"passing error", &entities.Employee{Name: "jason", Dob: "12-06-1998", City: "Bangalore",
			Majors: "MBA", DId: 1}, 1,
			&entities.Employee{}, errors.New("error")},
		{"rowsAffected is zero", &entities.Employee{Name: "jason", Dob: "12-06-1998", City: "Bangalore",
			Majors: "MBA", DId: 1}, 0,
			&entities.Employee{}, nil},
	}

	var s Store

	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	s = New(db)

	for i, tc := range testcases {
		mock.ExpectExec("Insert into employee values").
			WithArgs(sqlmock.AnyArg(), tc.input.Name, tc.input.Dob, tc.input.City, tc.input.Majors, tc.input.DId).
			WillReturnResult(sqlmock.NewResult(1, tc.rowsAffec)).WillReturnError(tc.err)

		actualOutput, _ := s.Create(context.TODO(), tc.input)

		tc.expectedOutput.ID = actualOutput.ID
		if !reflect.DeepEqual(actualOutput, tc.expectedOutput) {
			t.Errorf("test case %v %s : Expected %v \nGot %v testcase", i+1, tc.desc, tc.expectedOutput, actualOutput)
		}
	}
}

func TestUpdateEmployee(t *testing.T) {
	testcases := []struct {
		desc           string
		id             uuid.UUID
		input          *entities.Employee
		rowsAffec      int64
		expectedOutput *entities.Employee
		err            error
	}{
		{desc: "valid input", id: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			input: &entities.Employee{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), Name: "jason",
				Dob: "12-06-1998", City: "Bangalore", Majors: "MBA", DId: 1}, rowsAffec: 1,
			expectedOutput: &entities.Employee{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Name: "jason", Dob: "12-06-1998", City: "Bangalore", Majors: "MBA", DId: 1}, err: nil},
		{desc: "Passing error", id: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			input: &entities.Employee{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Name: "jason", Dob: "12-06-1998", City: "Bangalore", Majors: "MBA", DId: 1},
			rowsAffec: 1, expectedOutput: &entities.Employee{}, err: errors.New("error")},
		{desc: "rowsAffected is zero", id: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			input: &entities.Employee{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Name: "jason", Dob: "12-06-1998", City: "Bangalore", Majors: "MBA", DId: 1}, rowsAffec: 0,
			expectedOutput: &entities.Employee{}, err: nil},
	}

	var s Store

	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println(err)
	}

	s = New(db)

	for i, tc := range testcases {
		mock.ExpectExec("update employee set ").
			WithArgs(tc.input.ID, tc.input.Name, tc.input.Dob, tc.input.City, tc.input.Majors, tc.input.DId, tc.id).
			WillReturnResult(sqlmock.NewResult(1, tc.rowsAffec)).WillReturnError(tc.err)

		actualOutput, _ := s.Update(context.TODO(), uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), tc.input)
		if !reflect.DeepEqual(actualOutput, tc.expectedOutput) {
			t.Errorf("test case %v %s : Expected %v \nGot %v ", i+1, tc.desc, tc.expectedOutput, actualOutput)
		}
	}
}

func TestDeleteEmployee(t *testing.T) {
	testcases := []struct {
		desc           string
		id             uuid.UUID
		rowsAffec      int64
		expectedOutput int
		err            error
	}{
		{"If id in db", uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), 1, 204, nil},
		{"Passing error", uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), 1, 400, errors.New("error")},
		{"rowsAffected is zero", uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), 0, 400, nil},
	}

	var s Store

	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println(err)
	}

	s = New(db)

	for i, tc := range testcases {
		mock.ExpectExec("Delete from employee where id=?").
			WithArgs(tc.id).WillReturnResult(sqlmock.NewResult(1, tc.rowsAffec)).WillReturnError(tc.err)

		actualOutput, _ := s.Delete(context.TODO(), uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))
		if actualOutput != tc.expectedOutput {
			t.Errorf("test case %v %s : Expected %v \nGot %v ", i+1, tc.desc, tc.expectedOutput, actualOutput)
		}
	}
}

func TestGetEmployee(t *testing.T) {
	testcases := []struct {
		desc           string
		id             uuid.UUID
		expectedOutput entities.EmployeeAndDepartment
		err            error
	}{
		{desc: "valid id", id: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			expectedOutput: entities.EmployeeAndDepartment{ID: uuid.MustParse(
				"123e4567-e89b-12d3-a456-426614174000"), Name: "jason", Dob: "12-06-1998", City: "Bangalore",
				Majors: "CSE", Dept: entities.Department{ID: 2, Name: "TECH", FloorNo: 2}}, err: nil},
		{desc: "passing error", id: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			expectedOutput: entities.EmployeeAndDepartment{}, err: errors.New("error")},
	}

	var s Store

	db, mock, _ := sqlmock.New()

	s = New(db)

	for i, tc := range testcases {
		rows := mock.NewRows([]string{"id", "name", "dob", "city", "major", "deptid", "dept_name", "floor"}).
			AddRow(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), tc.expectedOutput.Name,
				tc.expectedOutput.Dob, tc.expectedOutput.City,
				tc.expectedOutput.Majors, tc.expectedOutput.Dept.ID, tc.expectedOutput.Dept.Name,
				tc.expectedOutput.Dept.FloorNo)
		mock.ExpectQuery("select").WithArgs(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")).
			WillReturnRows(rows).WillReturnError(tc.err)

		actualOutput, _ := s.Read(context.TODO(), uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))
		if !reflect.DeepEqual(actualOutput, tc.expectedOutput) {
			t.Errorf("test case %v %s : Expected %v \nGot %v ", i+1, tc.desc, tc.expectedOutput, actualOutput)
		}
	}
}

func TestGetAllWithCondition(t *testing.T) {
	testcases := []struct {
		desc              string
		name              string
		includeDepartment bool
		expectedOutput    []entities.EmployeeAndDepartment
		err               error
	}{
		{desc: "valid name and include true value", name: "jason", includeDepartment: true,
			expectedOutput: []entities.EmployeeAndDepartment{{
				ID:   uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Name: "jason", Dob: "12-06-1998", City: "Bangalore", Majors: "CSE", Dept: entities.Department{ID: 2,
					Name: "TECH", FloorNo: 2}}}, err: nil},
		{desc: "valid name and include false value", name: "jason",
			expectedOutput: []entities.EmployeeAndDepartment{{
				ID:   uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Name: "jason", Dob: "12-06-1998", City: "Bangalore", Majors: "CSE", Dept: entities.Department{ID: 2}}}},
		{"passing error", "jason", false,
			[]entities.EmployeeAndDepartment{{
				ID:   uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Name: "jason", Dob: "12-06-1998", City: "Bangalore", Majors: "CSE", Dept: entities.Department{ID: 2}}},
			errors.New("error")},
	}

	var s Store

	db, mock, _ := sqlmock.New()

	s = New(db)

	for i, tc := range testcases {
		row := mock.NewRows([]string{"id", "name", "dob", "city", "majors", "deptid", "deptName", "floorNo"}).
			AddRow(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), tc.expectedOutput[0].Name,
				tc.expectedOutput[0].Dob, tc.expectedOutput[0].City, tc.expectedOutput[0].Majors,
				tc.expectedOutput[0].Dept.ID, tc.expectedOutput[0].Dept.Name, tc.expectedOutput[0].Dept.FloorNo)
		mock.ExpectQuery("select").
			WithArgs(tc.name).WillReturnRows(row).WillReturnError(tc.err)

		actualOutput, err := s.ReadAll(store.Parameters{context.TODO(), tc.name,
			tc.includeDepartment})
		if err != tc.err {
			t.Errorf("test case %v %s : Expected %v \nGot %v ", i+1, tc.desc, tc.expectedOutput, actualOutput)
		}
	}
}

func TestGetAllWithOutCondition(t *testing.T) {
	testcases := []struct {
		desc              string
		name              string
		includeDepartment bool
		expectedOutput    []entities.EmployeeAndDepartment
	}{
		{desc: "no name and include true value", includeDepartment: true,
			expectedOutput: []entities.EmployeeAndDepartment{{
				ID:   uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Name: "jason", Dob: "12-06-1998", City: "Bangalore", Majors: "CSE", Dept: entities.Department{ID: 2,
					Name: "TECH", FloorNo: 2}}, {ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Name: "vk", Dob: "05-11-1997", City: "Bangalore", Majors: "MBA", Dept: entities.Department{ID: 1,
					Name: "HR", FloorNo: 1}}}},
	}

	var s Store

	db, mock, err := sqlmock.New()

	s = New(db)

	for i, tc := range testcases {
		rows := mock.NewRows([]string{"id", "name", "dob", "city", "major", "dept_id", "dept_name", "floor"}).
			AddRow(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), tc.expectedOutput[0].Name,
				tc.expectedOutput[0].Dob, tc.expectedOutput[0].City, tc.expectedOutput[0].Majors,
				tc.expectedOutput[0].Dept.ID, tc.expectedOutput[0].Dept.Name, tc.expectedOutput[0].Dept.FloorNo).
			AddRow(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), tc.expectedOutput[1].Name,
				tc.expectedOutput[1].Dob,
				tc.expectedOutput[1].City, tc.expectedOutput[1].Majors, tc.expectedOutput[1].Dept.ID,
				tc.expectedOutput[1].Dept.Name, tc.expectedOutput[1].Dept.FloorNo)
		mock.ExpectQuery("select").WillReturnRows(rows).WillReturnError(err)

		actualOutput, _ := s.ReadAll(store.Parameters{context.TODO(), tc.name,
			tc.includeDepartment})

		if !reflect.DeepEqual(actualOutput, tc.expectedOutput) {
			t.Errorf("test case %v %s : Expected %v \nGot %v ", i+1, tc.desc, tc.expectedOutput, actualOutput)
		}
	}
}

func TestGetWithNoRows(t *testing.T) {
	testcases := []struct {
		desc              string
		name              string
		includeDepartment bool
		expectedOutput    []entities.EmployeeAndDepartment
		err               error
	}{
		{"valid ", "jason", true, []entities.EmployeeAndDepartment{}, errorsHandler.NoData{Msg: "No Data"}},
	}
	var s Store

	db, mock, _ := sqlmock.New()

	s = New(db)
	for i, tc := range testcases {
		mock.ExpectQuery("select").
			WithArgs(tc.name).WillReturnRows(mock.NewRows([]string{})).WillReturnError(nil)

		actualOutput, err := s.ReadAll(store.Parameters{context.TODO(), tc.name, tc.includeDepartment})
		if err != tc.err {
			t.Errorf("test case %v %s : Expected %v \nGot %v ", i+1, tc.desc, tc.expectedOutput, actualOutput)
		}
	}
}

func TestGetDepartment(t *testing.T) {
	testcases := []struct {
		desc           string
		input          int
		expectedOutput entities.Department
		err            error
	}{
		{desc: "id in database", input: 1, expectedOutput: entities.Department{ID: 1, Name: "jason", FloorNo: 1}, err: nil},
		{"id not in database", 6, entities.Department{}, errors.New("err")},
		{"db error", 6, entities.Department{}, errors.New("err")},
	}

	var s Store

	db, mock, _ := sqlmock.New()

	s.Db = db

	for i, tc := range testcases {
		row := mock.NewRows([]string{"id", "name", "floorNo"}).AddRow(tc.expectedOutput.ID, tc.expectedOutput.Name, tc.expectedOutput.FloorNo)
		mock.ExpectQuery("select").WithArgs(tc.input).WillReturnRows(row).WillReturnError(tc.err)

		actualOutput, _ := s.ReadDepartment(context.TODO(), tc.input)
		if !reflect.DeepEqual(actualOutput, tc.expectedOutput) {
			t.Errorf("test case %v %s : Expected %v \nGot %v ", i+1, tc.desc, tc.expectedOutput, actualOutput)
		}
	}
}
