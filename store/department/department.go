package department

import (
	"EmployeeDepartment/entities"
	"database/sql"
)

type Store struct {
	Db *sql.DB
}

func New(db *sql.DB) Store {
	return Store{Db: db}
}

func (s Store) createDepartment(dept entities.Department) bool {
	res, err := s.Db.Exec("Insert into department values(?,?,?)", dept.Id, dept.Name, dept.FloorNo)
	if err != nil {
		return false
	}
	rowsAffected, err := res.RowsAffected()
	if rowsAffected == 1 {
		return true
	}
	return false
}

func (s Store) updateDepartmet(id int, dept entities.Department) bool {
	res, err := s.Db.Exec("Update department set id=id name=name floorNo=floorNo where id=?", dept.Id, dept.Name, dept.FloorNo, id)
	if err != nil {
		return false
	}
	rowsAffected, err := res.RowsAffected()
	if rowsAffected == 1 {
		return true
	}
	return false
}

func (s Store) deleteDepartment(id int) bool {
	res, err := s.Db.Exec("Delete from department where id=?", id)
	if err != nil {
		return false
	}
	rowsAffected, err := res.RowsAffected()
	if rowsAffected == 1 {
		return true
	}
	return false
}
