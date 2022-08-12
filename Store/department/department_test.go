package department

import (
	"EmployeeDepartment/Handler/Entities"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
)

func TestCreateDepartment(t *testing.T) {
	testCases := []struct {
		desc           string
		input          Entities.Department
		expectedOutput bool
	}{
		{"Valid input", Entities.Department{1, "HR", 1}, true},
		{"Invalid name", Entities.Department{2, "", 2}, true},
		{"Invalid floorNo", Entities.Department{2, "TECH", 0}, true},
	}
	var s Store

	db, mock, err := sqlmock.New()
	defer db.Close()

	s = New(db)
	for i, tc := range testCases {
		mock.ExpectExec("Insert into department values").
			WithArgs(tc.input.Id, tc.input.Name, tc.input.FloorNo).
			WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(err)
		actualOutput := s.createDepartment(tc.input)

		if actualOutput != tc.expectedOutput {
			t.Errorf("test case %v %s : Expected %v \nGot %v testcase", i+1, tc.desc, tc.expectedOutput, actualOutput)
		}
	}
}
