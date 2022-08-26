package department

import (
	"context"
	"database/sql"
	"net/http"

	"EmployeeDepartment/entities"
	"EmployeeDepartment/errorsHandler"
)

type Store struct {
	Db *sql.DB
}

func New(db *sql.DB) Store {
	return Store{Db: db}
}

func (s Store) Create(ctx context.Context, department entities.Department) (entities.Department, error) {
	res, err := s.Db.Exec("Insert into department values(?,?,?)", department.ID, department.Name, department.FloorNo)
	if err != nil {
		return entities.Department{}, err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 1 {
		return department, nil
	}

	return entities.Department{}, errorsHandler.AlreadyExists{Msg: "Already Exists"}
}

func (s Store) Update(ctx context.Context, id int, department entities.Department) (entities.Department, error) {
	res, err := s.Db.Exec("Update department set id=? ,name=? ,floor=? where id=?", department.ID, department.Name, department.FloorNo, id)
	if err != nil {
		return entities.Department{}, err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 1 {
		return department, nil
	}

	return entities.Department{}, &errorsHandler.AlreadyExists{Msg: "Already Exists"}
}

func (s Store) Delete(ctx context.Context, id int) (int, error) {
	res, err := s.Db.Exec("Delete from department where id=?", id)
	if err != nil {
		return http.StatusBadRequest, err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 1 {
		return http.StatusNoContent, nil
	}

	return http.StatusBadRequest, &errorsHandler.IDNotFound{Msg: "ID not found"}
}

func (s Store) GetDepartment(ctx context.Context, id int) (entities.Department, error) {
	var dept entities.Department
	res, err := s.Db.QueryContext(ctx, "select *  from department where id=?", id)
	if err != nil {
		return entities.Department{}, err
	}

	res.Next()
	res.Scan(&dept.ID, &dept.Name, &dept.FloorNo)
	return dept, nil

}
