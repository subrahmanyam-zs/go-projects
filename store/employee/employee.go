package employee

import (
	entities2 "EmployeeDepartment/entities"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

type Store struct {
	Db *sql.DB
}

func New(db *sql.DB) Store {
	return Store{Db: db}
}

func (s Store) createEmployee(emp entities2.Employee) bool {

	res, err := s.Db.Exec("Insert into employee values(?,?,?,?,?,?)", emp.Id, emp.Name, emp.Dob, emp.City, emp.Majors, emp.DId)
	if err != nil {
		fmt.Println(err)
		return false
	}
	rowsAffected, err := res.RowsAffected()
	if rowsAffected == 1 {
		return true
	}

	return false
}
func (s Store) updateEmployee(id uuid.UUID, emp entities2.Employee) bool {
	res, err := s.Db.Exec("Update employee set Id=id,Name=name,Dob=dob,City=city,Majors=majors,Did=did where Id=id", emp.Id, emp.Name, emp.Dob, emp.City, emp.Majors, emp.DId, id)
	if err != nil {
		return false
	}
	rowsAffected, err := res.RowsAffected()
	if rowsAffected == 1 {
		return true
	}
	return false
}

func (s Store) deleteEmployee(id uuid.UUID) bool {
	res, err := s.Db.Exec("Delete from employee where id=?", id)
	if err != nil {
		return false
	}
	rowsAffected, err := res.RowsAffected()
	if rowsAffected == 1 {
		return true
	}
	return false
}

func (s Store) getById(id uuid.UUID) (e entities2.EmployeeAndDepartment) {
	rows, err := s.Db.Query("select * from employee as e INNER JOIN department as d on e.Id=d.id where id=?", id.String())
	rows.Next()
	rows.Scan(&e.Id, &e.Name, &e.Dob, &e.City, &e.Majors, &e.Dept.Id, &e.Dept.Name, &e.Dept.FloorNo)
	rows.Close()
	if err != nil {
		return e
	}
	return
}

func (s Store) getAll(name string, includeDepartment bool) []entities2.EmployeeAndDepartment {
	index := 0
	if name != "" && includeDepartment {
		rows, err := s.Db.Query("select * from employee as e INNER JOIN department as d on e.Id=d.Id where name=name;", name)
		if err != nil {
			return []entities2.EmployeeAndDepartment{}
		}
		out := make([]entities2.EmployeeAndDepartment, 1, 1)
		for rows.Next() {
			rows.Scan(&out[index].Id, &out[index].Name, &out[index].Dob, &out[index].City, &out[index].Majors, &out[index].Dept.Id, &out[index].Dept.Name, &out[index].Dept.FloorNo)
			index++
		}
		return out
	} else if name != "" && !includeDepartment {
		rows, err := s.Db.Query("select * from employee as e INNER JOIN department as d on e.Id=d.Id where name=name;", name)
		if err != nil {
			return []entities2.EmployeeAndDepartment{}
		}
		out := make([]entities2.EmployeeAndDepartment, 1, 1)
		for rows.Next() {
			rows.Scan(&out[index].Id, &out[index].Name, &out[index].Dob, &out[index].City, &out[index].Majors, &out[index].Dept.Id, &out[index].Dept.Name, &out[index].Dept.FloorNo)
			index++
		}
		out[0].Dept.Name = ""
		out[0].Dept.FloorNo = 0
		return out
	} else {
		rows, err := s.Db.Query("select e.id,e.name,e.dob,e.city,e.major,d.id,d.floor from employee as e  INNER JOIN department as d on e.dept_id=d.id;")
		if err != nil {
			return []entities2.EmployeeAndDepartment{}
		}
		data := make([]entities2.EmployeeAndDepartment, 2, 2)
		for rows.Next() {

			rows.Scan(&data[index].Id, &data[index].Name, &data[index].Dob, &data[index].City, &data[index].Majors, &data[index].Dept.Id, &data[index].Dept.Name, &data[index].Dept.FloorNo)
			index++
		}

		return data
	}
}

func (s Store) getDepartment(id int) entities2.Department {
	var out entities2.Department
	rows, err := s.Db.Query("select * from department where id=Id", id)
	if err != nil {
		return entities2.Department{}
	}
	rows.Next()
	rows.Scan(&out.Id, &out.Name, &out.FloorNo)
	rows.Close()
	return out
}
