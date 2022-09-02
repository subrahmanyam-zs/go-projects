package employee

import (
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"

	"developer.zopsmart.com/go/gofr/Emp-Dept/entities"
	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func TestPost(t *testing.T) {
	testcases := []struct {
		desc           string
		input          entities.Employee
		rowsAffec      int64
		expectedOutput entities.Employee
		err            error
	}{
		{desc: "valid input", input: entities.Employee{Name: "jason", Dob: "12-06-1998", City: "Bangalore",
			Majors: "MBA", DeptID: 1}, rowsAffec: 1,
			expectedOutput: entities.Employee{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Name: "jason", Dob: "12-06-1998", City: "Bangalore", Majors: "MBA", DeptID: 1}},
		{"passing error", entities.Employee{Name: "jason", Dob: "12-06-1998", City: "Bangalore",
			Majors: "MBA", DeptID: 1}, 1,
			entities.Employee{}, errors.DB{}},
		{"rowsAffected is zero", entities.Employee{Name: "jason", Dob: "12-06-1998", City: "Bangalore",
			Majors: "MBA", DeptID: 1}, 0,
			entities.Employee{}, nil},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	ctx := gofr.NewContext(nil, nil, gofr.New())
	ctx.DataStore = datastore.DataStore{ORM: db}

	for i, tc := range testcases {
		mock.ExpectExec("Insert into employee values").
			WithArgs(sqlmock.AnyArg(), tc.input.Name, tc.input.Dob, tc.input.City, tc.input.Majors, tc.input.DeptID).
			WillReturnResult(sqlmock.NewResult(1, tc.rowsAffec)).WillReturnError(tc.err)

		actualOutput, err := New().Post(ctx, tc.input)

		if i == 2 {
			break
		}

		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("test case %v %s : Expected %v \nGot %v testcase", i+1, tc.desc, tc.expectedOutput, actualOutput)
		}
	}
}

func TestUpdateEmployee(t *testing.T) {
	testcases := []struct {
		desc           string
		id             uuid.UUID
		input          entities.Employee
		rowsAffec      int64
		expectedOutput entities.Employee
		err            error
	}{
		{desc: "valid input", id: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			input: entities.Employee{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), Name: "jason",
				Dob: "12-06-1998", City: "Bangalore", Majors: "MBA", DeptID: 1}, rowsAffec: 1,
			expectedOutput: entities.Employee{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Name: "jason", Dob: "12-06-1998", City: "Bangalore", Majors: "MBA", DeptID: 1}, err: nil},
		{desc: "Passing error", id: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			input: entities.Employee{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Name: "jason", Dob: "12-06-1998", City: "Bangalore", Majors: "MBA", DeptID: 1},
			rowsAffec: 1, expectedOutput: entities.Employee{}, err: errors.DB{}},
		{desc: "rowsAffected is zero", id: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			input: entities.Employee{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Name: "jason", Dob: "12-06-1998", City: "Bangalore", Majors: "MBA", DeptID: 1}, rowsAffec: 0,
			expectedOutput: entities.Employee{}, err: nil},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println(err)
	}

	defer db.Close()

	ctx := gofr.NewContext(nil, nil, gofr.New())
	ctx.DataStore = datastore.DataStore{ORM: db}

	for i, tc := range testcases {
		mock.ExpectExec("update employee set ").
			WithArgs(tc.id, tc.input.Name, tc.input.Dob, tc.input.City, tc.input.Majors, tc.input.DeptID, tc.id).
			WillReturnResult(sqlmock.NewResult(1, tc.rowsAffec)).WillReturnError(tc.err)

		actualOutput, err := New().Put(ctx, tc.id, tc.input)
		if !reflect.DeepEqual(actualOutput, tc.expectedOutput) {
			t.Errorf("test case %v %s : Expected %v error %v \nGot %v error %v ", i+1, tc.desc, tc.expectedOutput, tc.err, actualOutput, err)
		}
	}
}

func TestDelete(t *testing.T) {
	testcases := []struct {
		desc           string
		id             uuid.UUID
		rowsAffec      int64
		expectedOutput int
		err            error
	}{
		{"If id in db", uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), 1, 204, nil},
		{"Passing error", uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), 1, 400, errors.DB{}},
		{"rowsAffected is zero", uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), 0, 400, nil},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println(err)
	}

	defer db.Close()

	ctx := gofr.NewContext(nil, nil, gofr.New())
	ctx.DataStore = datastore.DataStore{ORM: db}

	for i, tc := range testcases {
		mock.ExpectExec("Delete from employee where id=?").
			WithArgs(tc.id).WillReturnResult(sqlmock.NewResult(1, tc.rowsAffec)).WillReturnError(tc.err)

		actualOutput, _ := New().Delete(ctx, tc.id)
		if actualOutput != tc.expectedOutput {
			t.Errorf("test case %v %s : Expected %v \nGot %v ", i+1, tc.desc, tc.expectedOutput, actualOutput)
		}
	}
}

func TestGet(t *testing.T) {
	testcases := []struct {
		desc           string
		id             uuid.UUID
		expectedOutput entities.EmpDept
		err            error
	}{
		{desc: "valid id", id: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			expectedOutput: entities.EmpDept{uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), "jason", "12-06-1998", "Bangalore", "CSE",
				entities.Department{2, "TECH", 2}}, err: nil},
		{desc: "passing error", id: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			expectedOutput: entities.EmpDept{}, err: errors.DB{}},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Println(err)
	}

	defer db.Close()

	ctx := gofr.NewContext(nil, nil, gofr.New())

	ctx.DataStore = datastore.DataStore{ORM: db}

	for i, tc := range testcases {
		rows := mock.NewRows([]string{"id", "name", "dob", "city", "major", "deptid", "dept_name", "floor"}).
			AddRow(tc.expectedOutput.ID, tc.expectedOutput.Name, tc.expectedOutput.Dob, tc.expectedOutput.City,
				tc.expectedOutput.Majors, tc.expectedOutput.Department.DeptID, tc.expectedOutput.Department.DeptName,
				tc.expectedOutput.Department.FloorNo)
		mock.ExpectQuery("select").WithArgs(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")).
			WillReturnRows(rows).WillReturnError(tc.err)

		actualOutput, err := New().Get(ctx, tc.id)
		if !reflect.DeepEqual(actualOutput, tc.expectedOutput) {
			t.Errorf("test case %v %s : Expected %v error %v \nGot %v error %v ", i+1, tc.desc, tc.expectedOutput, tc.err, actualOutput, err)
		}
	}
}

func TestGetAllWithCondition(t *testing.T) {
	testcases := []struct {
		desc              string
		name              string
		includeDepartment bool
		expectedOutput    []entities.EmpDept
		err               error
	}{
		{desc: "valid name and include true value", name: "jason", includeDepartment: true,
			expectedOutput: []entities.EmpDept{{
				uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				"jason",
				"12-06-1998",
				"Bangalore",
				"CSE",
				entities.Department{2, "TECH", 2},
			}}, err: nil},
		{desc: "valid name and include false value", name: "jason",
			expectedOutput: []entities.EmpDept{{
				uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				"jason",
				"12-06-1998",
				"Bangalore",
				"CSE",
				entities.Department{},
			}}},
		{"passing error", "jason", false,
			[]entities.EmpDept{{
				uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				"jason",
				"12-06-1998",
				"Bangalore",
				"CSE",
				entities.Department{},
			}},
			errors.DB{}},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	ctx := gofr.NewContext(nil, nil, gofr.New())

	ctx.DataStore = datastore.DataStore{ORM: db}

	for i, tc := range testcases {
		row := mock.NewRows([]string{"id", "name", "dob", "city", "majors", "deptid", "deptName", "floorNo"}).
			AddRow(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), tc.expectedOutput[0].Name,
				tc.expectedOutput[0].Dob, tc.expectedOutput[0].City, tc.expectedOutput[0].Majors,
				tc.expectedOutput[0].Department.DeptID, tc.expectedOutput[0].Department.DeptName, tc.expectedOutput[0].Department.FloorNo)
		mock.ExpectQuery("select").
			WithArgs(tc.name).WillReturnRows(row).WillReturnError(tc.err)

		actualOutput, err := New().GetAll(ctx, tc.name, tc.includeDepartment)
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
		expectedOutput    []entities.EmpDept
	}{
		{desc: "no name and include true value", includeDepartment: true,
			expectedOutput: []entities.EmpDept{{
				ID:   uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Name: "jason", Dob: "12-06-1998", City: "Bangalore", Majors: "CSE", Department: entities.Department{DeptID: 2,
					DeptName: "TECH", FloorNo: 2}}, {ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Name: "vk", Dob: "05-11-1997", City: "Bangalore", Majors: "MBA", Department: entities.Department{DeptID: 1,
					DeptName: "HR", FloorNo: 1}}}},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Println(err)
	}

	defer db.Close()

	ctx := gofr.NewContext(nil, nil, gofr.New())
	ctx.DataStore = datastore.DataStore{ORM: db}

	for i, tc := range testcases {
		rows := mock.NewRows([]string{"id", "name", "dob", "city", "major", "dept_id", "dept_name", "floor"}).
			AddRow(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), tc.expectedOutput[0].Name,
				tc.expectedOutput[0].Dob, tc.expectedOutput[0].City, tc.expectedOutput[0].Majors,
				tc.expectedOutput[0].Department.DeptID, tc.expectedOutput[0].Department.DeptName, tc.expectedOutput[0].Department.FloorNo).
			AddRow(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), tc.expectedOutput[1].Name,
				tc.expectedOutput[1].Dob,
				tc.expectedOutput[1].City, tc.expectedOutput[1].Majors, tc.expectedOutput[1].Department.DeptID,
				tc.expectedOutput[1].Department.DeptName, tc.expectedOutput[1].Department.FloorNo)
		mock.ExpectQuery("select").WillReturnRows(rows).WillReturnError(err)

		actualOutput, _ := New().GetAll(ctx, tc.name,
			tc.includeDepartment)

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
		expectedOutput    []entities.EmpDept
		err               error
	}{
		{"valid ", "jason", true, []entities.EmpDept{}, errors.DB{Err: fmt.Errorf("no data")}},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Println(err)
	}

	defer db.Close()

	ctx := gofr.NewContext(nil, nil, gofr.New())
	ctx.DataStore = datastore.DataStore{ORM: db}

	for i, tc := range testcases {
		mock.ExpectQuery("select").
			WithArgs(tc.name).WillReturnRows(mock.NewRows([]string{})).WillReturnError(nil)

		actualOutput, err := New().GetAll(ctx, tc.name, tc.includeDepartment)
		if !reflect.DeepEqual(actualOutput, tc.expectedOutput) {
			t.Errorf("test case %v %s : Expected %v error %v \nGot %v error %v ", i+1, tc.desc, tc.expectedOutput, tc.err, actualOutput, err)
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
		{desc: "id in database", input: 1, expectedOutput: entities.Department{DeptID: 1, DeptName: "jason", FloorNo: 1}, err: nil},
		{"id not in database", 6, entities.Department{}, errors.DB{}},
		{"db error", 6, entities.Department{}, errors.DB{}},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Println(err)
	}

	defer db.Close()

	ctx := gofr.NewContext(nil, nil, gofr.New())
	ctx.DataStore = datastore.DataStore{ORM: db}

	for i, tc := range testcases {
		row := mock.NewRows([]string{"id", "name", "floorNo"}).
			AddRow(tc.expectedOutput.DeptID, tc.expectedOutput.DeptName, tc.expectedOutput.FloorNo)
		mock.ExpectQuery("select").WithArgs(tc.input).WillReturnRows(row).WillReturnError(tc.err)

		actualOutput, _ := New().GetDepartment(ctx, tc.input)
		if !reflect.DeepEqual(actualOutput, tc.expectedOutput) {
			t.Errorf("test case %v %s : Expected %v \nGot %v ", i+1, tc.desc, tc.expectedOutput, actualOutput)
		}
	}
}
