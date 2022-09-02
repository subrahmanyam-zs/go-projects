package department

import (
	"fmt"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"developer.zopsmart.com/go/gofr/Emp-Dept/entities"
	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func TestPost(t *testing.T) {
	testCases := []struct {
		desc           string
		input          entities.Department
		rowsAffec      int64
		expectedOutput entities.Department
		err            error
	}{
		{desc: "Valid input", input: entities.Department{DeptID: 1, DeptName: "HR", FloorNo: 1}, rowsAffec: 1,
			expectedOutput: entities.Department{DeptID: 1, DeptName: "HR", FloorNo: 1}, err: nil},
		{desc: "passing error", input: entities.Department{DeptID: 1, DeptName: "HR", FloorNo: 1}, rowsAffec: 1, err: errors.DB{}},
		{desc: "Rows affected is zero", input: entities.Department{DeptID: 1, DeptName: "HR", FloorNo: 1}},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Println(err)
	}

	defer db.Close()

	ctx := gofr.NewContext(nil, nil, gofr.New())

	ctx.DataStore = datastore.DataStore{ORM: db}

	for i, tc := range testCases {
		mock.ExpectExec("Insert into department values").
			WithArgs(tc.input.DeptID, tc.input.DeptName, tc.input.FloorNo).
			WillReturnResult(sqlmock.NewResult(1, tc.rowsAffec)).WillReturnError(tc.err)

		actualOutput, _ := New().Post(ctx, tc.input)
		if actualOutput != tc.expectedOutput {
			t.Errorf("test case %v %s : Expected %v \nGot %v testcase", i+1, tc.desc, tc.expectedOutput, actualOutput)
		}
	}
}

func TestPutDepartment(t *testing.T) {
	testcases := []struct {
		desc           string
		id             int
		dataToUpdate   entities.Department
		rowsAffec      int64
		expectedOutput entities.Department
		err            error
	}{
		{desc: "valid input", id: 1, dataToUpdate: entities.Department{DeptID: 1, DeptName: "HR", FloorNo: 1}, rowsAffec: 1,
			expectedOutput: entities.Department{DeptID: 1, DeptName: "HR", FloorNo: 1}},
		{desc: "passing error", id: 1, dataToUpdate: entities.Department{DeptID: 1, DeptName: "HR", FloorNo: 1}, rowsAffec: 1,
			err: errors.DB{}},
		{desc: "rowsAffected is zero", id: 1, dataToUpdate: entities.Department{DeptID: 1, DeptName: "HR", FloorNo: 1}},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	ctx := gofr.NewContext(nil, nil, gofr.New())
	ctx.DataStore = datastore.DataStore{ORM: db}

	for i, tc := range testcases {
		mock.ExpectExec("Update").
			WithArgs(tc.dataToUpdate.DeptID, tc.dataToUpdate.DeptName, tc.dataToUpdate.FloorNo, tc.id).
			WillReturnResult(sqlmock.NewResult(1, tc.rowsAffec)).WillReturnError(tc.err)

		actualOutput, _ := New().Put(ctx, tc.id, tc.dataToUpdate)
		if actualOutput != tc.expectedOutput {
			t.Errorf("test case %v %s : Expected %v \nGot %v testcase", i+1, tc.desc, tc.expectedOutput, actualOutput)
		}
	}
}

func TestDeleteDepartment(t *testing.T) {
	testcases := []struct {
		desc           string
		id             int
		rowsAffec      int64
		expectedOutput int
		err            error
	}{
		{"valid id", 1, 1, 204, nil},
		{"passing error", 1, 1, 400, errors.DB{}},
		{"rowsAffected is zero", 1, 0, 400, nil},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	ctx := gofr.NewContext(nil, nil, gofr.New())
	ctx.DataStore = datastore.DataStore{ORM: db}

	for i, tc := range testcases {
		mock.ExpectExec("Delete from department where id=?").
			WithArgs(tc.id).WillReturnResult(sqlmock.NewResult(1, tc.rowsAffec)).WillReturnError(tc.err)

		actualOutput, _ := New().Delete(ctx, tc.id)
		if actualOutput != tc.expectedOutput {
			t.Errorf("test case %v %s : Expected %v \nGot %v testcase", i+1, tc.desc, tc.expectedOutput, actualOutput)
		}
	}
}
