package employee

import (
	entities2 "EmployeeDepartment/entities"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"net/http"
)

type Store struct {
	Db *sql.DB
}

func New(db *sql.DB) Store {
	return Store{Db: db}
}

func (s Store) Create(employee entities2.Employee) (entities2.Employee, error) {
	employee.Id = uuid.New()
	res, err := s.Db.Exec("Insert into employee values(?,?,?,?,?,?)", employee.Id, employee.Name, employee.Dob, employee.City, employee.Majors, employee.DId)
	if err != nil {
		return entities2.Employee{}, err
	}
	rowsAffected, err := res.RowsAffected()
	if rowsAffected == 1 {
		return employee, nil
	}
	return entities2.Employee{}, err
}

func (s Store) Update(id uuid.UUID, emp entities2.Employee) (entities2.Employee, error) {
	res, err := s.Db.Exec("Update employee set Id=?,Name=?,Dob=?,City=?,Majors=?,Did=? where Id=?", emp.Id, emp.Name, emp.Dob, emp.City, emp.Majors, emp.DId, id)
	if err != nil {
		return entities2.Employee{}, err
	}
	rowsAffected, err := res.RowsAffected()
	if rowsAffected == 1 {
		emp.Id = id
		return emp, nil
	}
	return entities2.Employee{}, errors.New("error")
}

func (s Store) Delete(id uuid.UUID) (int, error) {
	res, err := s.Db.Exec("Delete from employee where id=?", id)
	if err != nil {
		return http.StatusBadRequest, err
	}
	rowsAffected, err := res.RowsAffected()
	if rowsAffected == 1 {
		return http.StatusNoContent, nil
	}
	return http.StatusBadRequest, errors.New("error")
}

func (s Store) Read(id uuid.UUID) (entities2.EmployeeAndDepartment, error) {
	var out entities2.EmployeeAndDepartment
	rows := s.Db.QueryRow("select e.id,e.name,e.dob,e.city,e.majors,e.DId,d.name,d.floor from employee as e INNER JOIN department as d on e.DId=d.id where e.id=?", id)
	err := rows.Scan(&out.Id, &out.Name, &out.Dob, &out.City, &out.Majors, &out.Dept.Id, &out.Dept.Name, &out.Dept.FloorNo)
	if err != nil {
		return entities2.EmployeeAndDepartment{}, err
	}
	return out, nil
}

func (s Store) ReadAll(name string, includeDepartment bool) ([]entities2.EmployeeAndDepartment, error) {
	if name != "" && includeDepartment {
		rows, err := s.Db.Query("select e.id,e.name,e.dob,e.city,e.majors,d.id,d.name,d.floor from employee as e  INNER JOIN department as d on e.DId=d.Id where e.name=?;", name)
		if err != nil {
			return []entities2.EmployeeAndDepartment{}, err
		}
		data := make([]entities2.EmployeeAndDepartment, 0, 0)
		for rows.Next() {
			var temp entities2.EmployeeAndDepartment
			rows.Scan(&temp.Id, &temp.Name, &temp.Dob, &temp.City, &temp.Majors, &temp.Dept.Id, &temp.Dept.Name, &temp.Dept.FloorNo)
			data = append(data, temp)
		}
		return data, nil
	} else if name != "" && !includeDepartment {
		rows, err := s.Db.Query("select e.id,e.name,e.dob,e.city,e.majors,d.id,d.name,d.floor from employee as e INNER JOIN department as d on e.DId=d.Id where e.name=?;", name)
		if err != nil {
			return []entities2.EmployeeAndDepartment{}, err
		}
		data := make([]entities2.EmployeeAndDepartment, 0, 0)
		for rows.Next() {
			var temp entities2.EmployeeAndDepartment
			rows.Scan(&temp.Id, &temp.Name, &temp.Dob, &temp.City, &temp.Majors, &temp.Dept.Id, &temp.Dept.Name, &temp.Dept.FloorNo)
			temp.Dept.FloorNo = 0
			temp.Dept.Name = ""
			data = append(data, temp)
		}
		return data, nil
	} else {
		rows, err := s.Db.Query("select e.id,e.name,e.dob,e.city,e.majors,d.id,d.name,d.floor from employee as e  INNER JOIN department as d on e.DId=d.id;")
		if err != nil {
			return []entities2.EmployeeAndDepartment{}, err
		}
		data := make([]entities2.EmployeeAndDepartment, 0, 0)
		for rows.Next() {
			var temp entities2.EmployeeAndDepartment
			rows.Scan(&temp.Id, &temp.Name, &temp.Dob, &temp.City, &temp.Majors, &temp.Dept.Id, &temp.Dept.Name, &temp.Dept.FloorNo)
			data = append(data, temp)
		}
		return data, nil
	}
}

func (s Store) ReadDepartment(id int) (entities2.Department, error) {
	var out entities2.Department
	rows, err := s.Db.Query("select * from department where id=?", id)
	if err != nil {
		return entities2.Department{}, err
	}
	rows.Next()
	rows.Scan(&out.Id, &out.Name, &out.FloorNo)
	rows.Close()

	return out, nil
}
