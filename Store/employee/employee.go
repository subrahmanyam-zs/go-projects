package employee

import (
	"EmployeeDepartment/Handler/Entities"
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

func (s Store) createEmployee(emp Entities.Employee) bool {

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
func (s Store) updateEmployee(id uuid.UUID, emp Entities.Employee) bool {
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

func (s Store) getById(id uuid.UUID) (e Entities.EmployeeAndDepartment) {
	rows, err := s.Db.Query("select * from employee as e INNER JOIN department as d on e.Id=d.id where id=?", id.String())
	rows.Next()
	rows.Scan(&e.Id, &e.Name, &e.Dob, &e.City, &e.Majors, &e.Dept.Id, &e.Dept.Name, &e.Dept.FloorNo)
	rows.Close()
	if err != nil {
		return e
	}
	return
}
