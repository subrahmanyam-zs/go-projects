package department

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"EmployeeDepartment/entities"
)

func TestCreateDepartment(t *testing.T) {
	testCases := []struct {
		desc           string
		input          entities.Department
		rowsAffec      int64
		expectedOutput entities.Department
		err            error
	}{
		{desc: "Valid input", input: entities.Department{ID: 1, Name: "HR", FloorNo: 1}, rowsAffec: 1,
			expectedOutput: entities.Department{ID: 1, Name: "HR", FloorNo: 1}},
		{desc: "passing error", input: entities.Department{ID: 1, Name: "HR", FloorNo: 1}, rowsAffec: 1, err: errors.New("err")},
		{desc: "Rows affected is zero", input: entities.Department{ID: 1, Name: "HR", FloorNo: 1}},
	}

	var s Store

	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println(err)
	}

	defer db.Close()

	s = New(db)

	for i, tc := range testCases {
		mock.ExpectExec("Insert into department values").
			WithArgs(tc.input.ID, tc.input.Name, tc.input.FloorNo).
			WillReturnResult(sqlmock.NewResult(1, tc.rowsAffec)).WillReturnError(tc.err)

		actualOutput, _ := s.Create(context.TODO(), tc.input)
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
		{desc: "valid input", id: 1, dataToUpdate: entities.Department{ID: 1, Name: "HR", FloorNo: 1}, rowsAffec: 1,
			expectedOutput: entities.Department{ID: 1, Name: "HR", FloorNo: 1}},
		{desc: "passing error", id: 1, dataToUpdate: entities.Department{ID: 1, Name: "HR", FloorNo: 1}, rowsAffec: 1,
			err: errors.New("error")},
		{desc: "rowsAffected is zero", id: 1, dataToUpdate: entities.Department{ID: 1, Name: "HR", FloorNo: 1}},
	}

	var s Store

	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println(err)
	}

	s = New(db)

	for i, tc := range testcases {
		mock.ExpectExec("Update").
			WithArgs(tc.dataToUpdate.ID, tc.dataToUpdate.Name, tc.dataToUpdate.FloorNo, tc.id).
			WillReturnResult(sqlmock.NewResult(1, tc.rowsAffec)).WillReturnError(tc.err)

		actualOutput, _ := s.Update(context.TODO(), tc.id, tc.dataToUpdate)
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
		{"passing error", 1, 1, 400, errors.New("error")},
		{"rowsAffected is zero", 1, 0, 400, nil},
	}

	var s Store

	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println(err)
	}

	s.Db = db

	for i, tc := range testcases {
		mock.ExpectExec("Delete from department where id=?").
			WithArgs(tc.id).WillReturnResult(sqlmock.NewResult(1, tc.rowsAffec)).WillReturnError(tc.err)

		actualOutput, _ := s.Delete(context.TODO(), tc.id)
		if actualOutput != tc.expectedOutput {
			t.Errorf("test case %v %s : Expected %v \nGot %v testcase", i+1, tc.desc, tc.expectedOutput, actualOutput)
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
		{"valid id", 1, entities.Department{1, "HR", 1}, nil},
		{"passing error", 1, entities.Department{}, errors.New("err")},
	}

	var s Store

	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println(err)
	}

	s.Db = db
	for i, tc := range testcases {
		rows := mock.NewRows([]string{"id", "string", "floor"}).AddRow(tc.expectedOutput.ID, tc.expectedOutput.Name, tc.expectedOutput.FloorNo)
		mock.ExpectQuery("select").WithArgs(tc.input).WillReturnRows(rows).WillReturnError(tc.err)

		actualOutput, _ := s.GetDepartment(context.TODO(), tc.input)
		if actualOutput != tc.expectedOutput {
			t.Errorf("test case %v %s : Expected %v \nGot %v testcase", i+1, tc.desc, tc.expectedOutput, actualOutput)
		}
	}
}