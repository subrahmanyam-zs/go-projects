package employee

import (
	"EmployeeDepartment/Handler/Entities"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"reflect"
	"testing"
)

var uid uuid.UUID = uuid.New()

func TestcreateEmployee(t *testing.T) {
	testcases := []struct {
		desc           string
		input          Entities.Employee
		expectedOutput bool
	}{
		{"Invalid input", Entities.Employee{uid, "jason", "12-06-1999", "Bangalore", "CSE", 2}, true},
		{"Valid input", Entities.Employee{uid, "jason", "12-06-1998", "Bangalore", "MBA", 1}, true},
	}

	var s Store

	db, mock, err := sqlmock.New()
	defer db.Close()

	s = New(db)
	for i, tc := range testcases {
		mock.ExpectExec("Insert into employee values").
			WithArgs(tc.input.Id, tc.input.Name, tc.input.Dob, tc.input.City, tc.input.Majors, tc.input.DId).
			WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(err)
		actualOutput := s.createEmployee(tc.input)

		if actualOutput != tc.expectedOutput {
			t.Errorf("test case %v %s : Expected %v \nGot %v testcase", i+1, tc.desc, tc.expectedOutput, actualOutput)
		}
	}
}

func TestUpdateEmployee(t *testing.T) {
	testcases := []struct {
		desc           string
		id             uuid.UUID
		input          Entities.Employee
		expectedOutput bool
	}{
		{"Invalid input", uid, Entities.Employee{uid, "jason", "12-06-1999", "Bangalore", "CSE", 2}, true},
		{"", uid, Entities.Employee{uid, "jason", "12-06-1998", "Bangalore", "MBA", 1}, true},
	}
	var s Store
	db, mock, err := sqlmock.New()
	s = New(db)

	for i, tc := range testcases {
		mock.ExpectExec("Update employee set Id=id,Name=name,Dob=dob,City=city,Majors=majors,Did=did where Id=id").
			WithArgs(tc.input.Id, tc.input.Name, tc.input.Dob, tc.input.City, tc.input.Majors, tc.input.DId, tc.id).
			WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(err)
		actualOutput := s.updateEmployee(tc.id, tc.input)
		if actualOutput != tc.expectedOutput {
			t.Errorf("test case %v %s : Expected %v \nGot %v ", i+1, tc.desc, tc.expectedOutput, actualOutput)
		}
	}
}

func TestDeleteEmployee(t *testing.T) {
	testcases := []struct {
		desc           string
		id             uuid.UUID
		expectedOutput bool
	}{
		{"If id in db", uid, true},
		{"If id not in db", uid, true},
	}

	var s Store

	db, mock, err := sqlmock.New()

	s = New(db)

	for i, tc := range testcases {
		mock.ExpectExec("Delete from employee where id=?").WithArgs(tc.id).WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(err)
		actualOutput := s.deleteEmployee(tc.id)
		if actualOutput != tc.expectedOutput {
			t.Errorf("test case %v %s : Expected %v \nGot %v ", i+1, tc.desc, tc.expectedOutput, actualOutput)
		}
	}
}

func TestGetEmployee(t *testing.T) {
	testcases := []struct {
		desc           string
		id             uuid.UUID
		expectedOutput Entities.EmployeeAndDepartment
	}{
		{"Valid id", uid, Entities.EmployeeAndDepartment{uid.String(), "jason", "12-06-1998", "Bangalore", "CSE", Entities.Department{2, "TECH", 2}}},
	}
	var s Store
	db, mock, err := sqlmock.New()
	s = New(db)

	for i, tc := range testcases {
		row := mock.NewRows([]string{"id", "name", "dob", "city", "majors", "deptid", "deptname", "floorNo"}).
			AddRow(tc.expectedOutput.Id, tc.expectedOutput.Name, tc.expectedOutput.Dob, tc.expectedOutput.City, tc.expectedOutput.Majors, tc.expectedOutput.Dept.Id, tc.expectedOutput.Dept.Name, tc.expectedOutput.Dept.FloorNo)
		mock.ExpectQuery("select (.?) from employee as e INNER JOIN department as d on e.Id=d.id where id=?").WithArgs(tc.id.String()).
			WillReturnRows(row).WillReturnError(err)
		actualOutput := s.getById(tc.id)

		if !reflect.DeepEqual(actualOutput, tc.expectedOutput) {
			t.Errorf("test case %v %s : Expected %v \nGot %v ", i+1, tc.desc, tc.expectedOutput, actualOutput)
		}
	}
}
