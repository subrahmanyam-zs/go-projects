package department

import (
	"net/http"

	"developer.zopsmart.com/go/gofr/Emp-Dept/entities"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Department struct {
}

func New() Department {
	return Department{}
}

func (d Department) Post(ctx *gofr.Context, dept entities.Department) (interface{}, error) {
	res, err := ctx.DB().Exec("Insert into department values(?,?,?)", dept.DeptID, dept.DeptName, dept.FloorNo)
	if err != nil {
		return entities.Department{}, err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 1 {
		return dept, nil
	}

	return entities.Department{}, errors.EntityAlreadyExists{}
}

func (d Department) Put(ctx *gofr.Context, id int, dept entities.Department) (interface{}, error) {
	res, err := ctx.DB().Exec("Update department set id=? ,name=? ,floor=? where id=?", dept.DeptID, dept.DeptName, dept.FloorNo, id)
	if err != nil {
		return entities.Department{}, err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 1 {
		return dept, nil
	}

	return entities.Department{}, errors.EntityNotFound{"department", "id"}
}

func (d Department) Delete(ctx *gofr.Context, id int) (int, error) {
	res, err := ctx.DB().Exec("Delete from department where id=?", id)
	if err != nil {
		return http.StatusBadRequest, err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 1 {
		return http.StatusNoContent, nil
	}

	return http.StatusBadRequest, errors.EntityNotFound{"department", "id"}
}
