package department

import (
	"EmployeeDepartment/entities"
	"database/sql"
	"errors"
	"fmt"
)

type Store struct {
	Db *sql.DB
}

func New(db *sql.DB) Store {
	return Store{Db: db}
}

func (s Store) Create(department entities.Department) (entities.Department, error) {
	res, err := s.Db.Exec("Insert into department values(?,?,?)", department.Id, department.Name, department.FloorNo)
	if err != nil {
		return entities.Department{}, err
	}
	rowsAffected, err := res.RowsAffected()
	fmt.Println(rowsAffected)
	if rowsAffected == 1 {
		return department, nil
	}
	return entities.Department{}, err
}

func (s Store) Update(id int, department entities.Department) (entities.Department, error) {
	res, err := s.Db.Exec("Update department set name=? ,floor=? where id=?", department.Name, department.FloorNo, id)
	if err != nil {
		return entities.Department{}, err
	}
	rowsAffected, err := res.RowsAffected()
	if rowsAffected == 1 {
		department.Id = id
		return department, nil
	}
	return entities.Department{}, errors.New("error")
}

func (s Store) Delete(id int) (int, error) {
	res, err := s.Db.Exec("Delete from department where id=?", id)
	if err != nil {
		return 400, err
	}
	rowsAffected, err := res.RowsAffected()
	if rowsAffected == 1 {
		return 204, nil
	}
	return 400, err
}
