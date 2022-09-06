package datastore

import (
	"developer.zopsmart.com/go/gofr/examples/sample-api/entity"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"fmt"
)

type employee struct{}

func New() employee {
	return employee{}
}

func (e employee) Post(ctx *gofr.Context, emp entity.Employee) (entity.Employee, error) {
	res, err := ctx.DB().Exec("Insert into emp values(?,?,?,?)", emp.ID, emp.Name, emp.City, emp.Majors)
	if err != nil {
		return entity.Employee{}, errors.DB{Err: fmt.Errorf("query error")}
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return entity.Employee{}, errors.DB{Err: fmt.Errorf("error in rowsAfected")}
	} else if rowsAffected == 1 {
		return emp, nil
	}

	return entity.Employee{}, errors.EntityAlreadyExists{}
}

func (e employee) Put(ctx *gofr.Context, id string, emp entity.Employee) (entity.Employee, error) {
	res, err := ctx.DB().Exec("Update emp set id=?, name=?, city=?, majors=? where id=?", emp.ID, emp.Name, emp.City, emp.Majors, id)
	if err != nil {
		return entity.Employee{}, errors.DB{Err: err}
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return entity.Employee{}, errors.DB{Err: fmt.Errorf("error in rowsAfected")}
	} else if rowsAffected == 1 {
		return emp, nil
	}

	return entity.Employee{}, errors.EntityAlreadyExists{}
}

func (e employee) Delete(ctx *gofr.Context, id string) (int, error) {
	res, err := ctx.DB().Exec("Delete from emp where id=?", id)
	if err != nil {
		fmt.Println(err)
		return 400, errors.DB{Err: err}
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 400, errors.DB{Err: err}
	} else if rowsAffected == 1 {
		return 204, nil
	}

	return 404, errors.MissingParam{Param: []string{"id not found"}}
}

func (e employee) Get(ctx *gofr.Context, id string) (entity.Employee, error) {
	var emp entity.Employee

	row := ctx.DB().QueryRow("select *from emp where id=?", id)

	err := row.Scan(&emp.ID, &emp.Name, &emp.City, &emp.Majors)
	if err != nil {
		return entity.Employee{}, errors.DB{Err: err}
	}

	return emp, nil
}

func (e employee) GetAll(ctx *gofr.Context) ([]entity.Employee, error) {
	var emp []entity.Employee

	rows, err := ctx.DB().Query("select *from emp")
	if err != nil {
		return []entity.Employee{}, err
	}

	for rows.Next() {
		var e entity.Employee
		err := rows.Scan(&e.ID, &e.Name, &e.City, &e.Majors)
		if err != nil {
			return emp, errors.DB{Err: err}
		}
		emp = append(emp, e)
	}

	return emp, nil
}
