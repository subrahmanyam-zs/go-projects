package employee

import (
	entities2 "EmployeeDepartment/entities"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"reflect"
	"testing"
)

var uid uuid.UUID = uuid.New()

func TestCreateEmployee(t *testing.T) {
	testcases := []struct {
		desc           string
		input          entities2.Employee
		expectedOutput entities2.Employee
	}{
		{"valid input", entities2.Employee{uid, "jason", "12-06-1998", "Bangalore", "MBA", 1}, entities2.Employee{uid, "jason", "12-06-1998", "Bangalore", "MBA", 1}},
	}
	var s Store
	db, mock, err := sqlmock.New()
	defer db.Close()
	s = New(db)
	for i, tc := range testcases {
		mock.ExpectExec("Insert into employee values").
			WithArgs(tc.input.Id, tc.input.Name, tc.input.Dob, tc.input.City, tc.input.Majors, tc.input.DId).
			WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(err)
		actualOutput, _ := s.Create(tc.input)
		if actualOutput != tc.expectedOutput {
			t.Errorf("test case %v %s : Expected %v \nGot %v testcase", i+1, tc.desc, tc.expectedOutput, actualOutput)
		}
	}
}

func TestUpdateEmployee(t *testing.T) {
	testcases := []struct {
		desc           string
		id             uuid.UUID
		input          entities2.Employee
		expectedOutput entities2.Employee
	}{
		{"valid input", uid, entities2.Employee{uid, "jason", "12-06-1998", "Bangalore", "MBA", 1}, entities2.Employee{uid, "jason", "12-06-1998", "Bangalore", "MBA", 1}},
	}
	var s Store
	db, mock, err := sqlmock.New()
	s = New(db)

	for i, tc := range testcases {
		mock.ExpectExec("Update employee set Id=id,Name=name,Dob=dob,City=city,Majors=majors,Did=did where Id=id").
			WithArgs(tc.input.Id, tc.input.Name, tc.input.Dob, tc.input.City, tc.input.Majors, tc.input.DId, tc.id).
			WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(err)
		actualOutput, _ := s.Update(tc.id, tc.input)
		if actualOutput != tc.expectedOutput {
			t.Errorf("test case %v %s : Expected %v \nGot %v ", i+1, tc.desc, tc.expectedOutput, actualOutput)
		}
	}
}

func TestDeleteEmployee(t *testing.T) {
	testcases := []struct {
		desc           string
		id             uuid.UUID
		expectedOutput int
	}{
		{"If id in db", uid, 204},
		{"If id not in db", uid, 204},
	}
	var s Store
	db, mock, err := sqlmock.New()
	s = New(db)
	for i, tc := range testcases {
		mock.ExpectExec("Delete from employee where id=?").WithArgs(tc.id).WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(err)
		actualOutput, _ := s.Delete(tc.id)
		if actualOutput != tc.expectedOutput {
			t.Errorf("test case %v %s : Expected %v \nGot %v ", i+1, tc.desc, tc.expectedOutput, actualOutput)
		}
	}
}

func TestGetEmployee(t *testing.T) {
	testcases := []struct {
		desc           string
		id             uuid.UUID
		expectedOutput entities2.EmployeeAndDepartment
	}{
		{"valid id", uid, entities2.EmployeeAndDepartment{uid.String(), "jason", "12-06-1998", "Bangalore", "CSE", entities2.Department{2, "TECH", 2}}},
	}
	var s Store
	db, mock, err := sqlmock.New()
	s = New(db)

	for i, tc := range testcases {
		row := mock.NewRows([]string{"id", "name", "dob", "city", "majors", "deptid", "deptname", "floorNo"}).
			AddRow(tc.expectedOutput.Id, tc.expectedOutput.Name, tc.expectedOutput.Dob, tc.expectedOutput.City, tc.expectedOutput.Majors, tc.expectedOutput.Dept.Id, tc.expectedOutput.Dept.Name, tc.expectedOutput.Dept.FloorNo)
		mock.ExpectQuery("select (.?) from employee as e INNER JOIN department as d on e.Id=d.id where id=?").WithArgs(tc.id.String()).
			WillReturnRows(row).WillReturnError(err)
		actualOutput, _ := s.Read(tc.id)
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
		expectedOutput    []entities2.EmployeeAndDepartment
	}{
		{"tcid name and include department as true tcue", "jason", true, []entities2.EmployeeAndDepartment{{uid.String(), "jason", "12-06-1998", "Bangalore", "CSE", entities2.Department{2, "TECH", 2}}}},
		{"tcid name and include department as false tcue", "jason", false, []entities2.EmployeeAndDepartment{{uid.String(), "jason", "12-06-1998", "Bangalore", "CSE", entities2.Department{2, "", 0}}}},
	}

	var s Store
	db, mock, err := sqlmock.New()
	s = New(db)
	for i, tc := range testcases {
		row := mock.NewRows([]string{"id", "name", "dob", "city", "majors", "deptid", "deptName", "floorNo"}).
			AddRow(tc.expectedOutput[0].Id, tc.expectedOutput[0].Name, tc.expectedOutput[0].Dob, tc.expectedOutput[0].City, tc.expectedOutput[0].Majors, tc.expectedOutput[0].Dept.Id, tc.expectedOutput[0].Dept.Name, tc.expectedOutput[0].Dept.FloorNo)
		mock.ExpectQuery("select (.?) from employee as e INNER JOIN department as d on e.Id=d.Id where name=name;").
			WithArgs(tc.name).WillReturnRows(row).WillReturnError(err)
		actualOutput, _ := s.ReadAll(tc.name, tc.includeDepartment)
		if !reflect.DeepEqual(actualOutput, tc.expectedOutput) {
			t.Errorf("test case %v %s : Expected %v \nGot %v ", i+1, tc.desc, tc.expectedOutput, actualOutput)
		}
	}
}

func TestGetAllWithOutCondition(t *testing.T) {
	testcases := []struct {
		desc              string
		name              string
		includeDepartment bool
		expectedOutput    []entities2.EmployeeAndDepartment
	}{
		{"no name and include department as true value", "", true, []entities2.EmployeeAndDepartment{{uid.String(), "jason", "12-06-1998", "Bangalore", "CSE", entities2.Department{2, "TECH", 2}}, {uid.String(), "vk", "05-11-1997", "Bangalore", "MBA", entities2.Department{1, "HR", 1}}}},
	}

	var s Store
	db, mock, err := sqlmock.New()
	s = New(db)
	for i, tc := range testcases {

		rows := mock.NewRows([]string{"id", "name", "dob", "city", "major", "dept_id", "dept_name", "floor"}).
			AddRow(tc.expectedOutput[0].Id, tc.expectedOutput[0].Name, tc.expectedOutput[0].Dob, tc.expectedOutput[0].City, tc.expectedOutput[0].Majors, tc.expectedOutput[0].Dept.Id, tc.expectedOutput[0].Dept.Name, tc.expectedOutput[0].Dept.FloorNo).
			AddRow(tc.expectedOutput[1].Id, tc.expectedOutput[1].Name, tc.expectedOutput[1].Dob, tc.expectedOutput[1].City, tc.expectedOutput[1].Majors, tc.expectedOutput[1].Dept.Id, tc.expectedOutput[1].Dept.Name, tc.expectedOutput[1].Dept.FloorNo)
		mock.ExpectQuery("select e.id,e.name,e.dob,e.city,e.major,d.id,d.floor from employee as e INNER JOIN department as d on e.dept_id=d.id;").WillReturnRows(rows).WillReturnError(err)
		actualOutput, _ := s.ReadAll(tc.name, tc.includeDepartment)
		if !reflect.DeepEqual(actualOutput, tc.expectedOutput) {
			t.Errorf("test case %v %s : Expected %v \nGot %v ", i+1, tc.desc, tc.expectedOutput, actualOutput)
		}
	}
}

func TestGetDepartment(t *testing.T) {
	testcases := []struct {
		desc           string
		input          int
		expectedOutput entities2.Department
	}{
		{"id in database", 1, entities2.Department{1, "jason", 1}},
		{"id not in database", 6, entities2.Department{}},
	}
	var s Store
	db, mock, err := sqlmock.New()
	s.Db = db
	for i, tc := range testcases {
		row := mock.NewRows([]string{"id", "name", "floorNo"}).AddRow(tc.expectedOutput.Id, tc.expectedOutput.Name, tc.expectedOutput.FloorNo)
		mock.ExpectQuery("select (.+) from department where id=Id").WithArgs(tc.input).WillReturnRows(row).WillReturnError(err)
		actualOutput, _ := s.ReadDepartment(tc.input)
		if !reflect.DeepEqual(actualOutput, tc.expectedOutput) {
			t.Errorf("test case %v %s : Expected %v \nGot %v ", i+1, tc.desc, tc.expectedOutput, actualOutput)
		}
	}
}
