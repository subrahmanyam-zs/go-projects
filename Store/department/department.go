package department

import (
	"EmployeeDepartment/Handler/Entities"
	"database/sql"
)

type Store struct {
	Db *sql.DB
}

func New(db *sql.DB) Store {
	return Store{Db: db}
}

func (s Store) createDepartment(dept Entities.Department) bool {
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
