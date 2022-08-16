package department

import (
	"EmployeeDepartment/entities"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
)

func TestCreateDepartment(t *testing.T) {
	testCases := []struct {
		desc           string
		input          entities.Department
		expectedOutput entities.Department
	}{
		{"Valid input", entities.Department{1, "HR", 1}, entities.Department{1, "HR", 1}},
		{"Invalid name", entities.Department{2, "", 2}, entities.Department{}},
		{"Invalid floorNo", entities.Department{2, "TECH", 0}, entities.Department{}},
	}
	var s Store

	db, mock, err := sqlmock.New()
	defer db.Close()

	s = New(db)
	for i, tc := range testCases {
		mock.ExpectExec("Insert into department values").
			WithArgs(tc.input.Id, tc.input.Name, tc.input.FloorNo).
			WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(err)
		actualOutput, _ := s.Create(tc.input)

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
		expectedOutput entities.Department
	}{
		{"valid input", 1, entities.Department{1, "HR", 1}, entities.Department{1, "HR", 1}},
		{"invalid id", 0, entities.Department{1, "HR", 1}, entities.Department{}},
		{"Invalid name", 2, entities.Department{2, "", 2}, entities.Department{}},
		{"Invalid floorNo", 2, entities.Department{2, "TECH", 0}, entities.Department{}},
	}
	var s Store
	db, mock, err := sqlmock.New()
	s = New(db)
	//s.Db = db
	for i, tc := range testcases {
		mock.ExpectExec("Update department set id=id name=name floorNo=floorNo where id=?").
			WithArgs(tc.dataToUpdate.Id, tc.dataToUpdate.Name, tc.dataToUpdate.FloorNo, tc.id).
			WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(err)
		actualOutput, _ := s.Update(tc.id, tc.dataToUpdate)
		if actualOutput != tc.expectedOutput {
			t.Errorf("test case %v %s : Expected %v \nGot %v testcase", i+1, tc.desc, tc.expectedOutput, actualOutput)
		}
	}
}

func TestDeleteDepartment(t *testing.T) {
	testcases := []struct {
		desc           string
		id             int
		expectedOutput int
	}{
		{"valid id", 1, 204},
		{"invalid id", 0, 204},
	}
	var s Store
	db, mock, err := sqlmock.New()
	s.Db = db

	for i, tc := range testcases {
		mock.ExpectExec("Delete from department where id=?").
			WithArgs(tc.id).WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(err)
		actualOutput, _ := s.Delete(tc.id)
		if actualOutput != tc.expectedOutput {
			t.Errorf("test case %v %s : Expected %v \nGot %v testcase", i+1, tc.desc, tc.expectedOutput, actualOutput)
		}
	}
}
