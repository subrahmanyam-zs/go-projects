package employee

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"developer.zopsmart.com/go/gofr/Emp-Dept/entities"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Store struct {
}

func New() Store {
	return Store{}
}

func (e Store) Post(ctx *gofr.Context, emp entities.Employee) (interface{}, error) {
	var uid uuid.UUID = uuid.New()

	res, err := ctx.DB().Exec("Insert into employee values(?,?,?,?,?,?)", uid, emp.Name, emp.Dob, emp.City, emp.Majors, emp.DeptID)
	if err != nil {
		return entities.Employee{}, err
	}

	rowsAffected, _ := res.RowsAffected()

	if rowsAffected == 1 {
		emp.ID = uid
		return emp, nil
	}

	return entities.Employee{}, errors.EntityAlreadyExists{}
}

func (e Store) Put(ctx *gofr.Context, id uuid.UUID, emp entities.Employee) (interface{}, error) {
	res, err := ctx.DB().Exec("update employee set id=?, name=?,dob=?,city=?,majors=?,dId=? "+
		"where id=?;", id, emp.Name, emp.Dob, emp.City, emp.Majors, emp.DeptID, id)
	if err != nil {
		return entities.Employee{}, err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 1 {
		emp.ID = id
		return emp, nil
	}

	return entities.Employee{}, errors.EntityAlreadyExists{}
}

func (e Store) Delete(ctx *gofr.Context, id uuid.UUID) (int, error) {
	res, err := ctx.DB().Exec("Delete from employee where id=?", id)
	if err != nil {
		return http.StatusBadRequest, err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 1 {
		return http.StatusNoContent, nil
	}

	return http.StatusBadRequest, errors.EntityNotFound{Entity: "employee", ID: "id"}
}

func (e Store) Get(ctx *gofr.Context, id uuid.UUID) (interface{}, error) {
	var out entities.EmpDept

	rows := ctx.DB().QueryRow(
		"select e.id,e.name,e.dob,e.city,e.majors,d.id,d.name,d.floor from employee"+
			" as e inner join department as d on e.DId=d.id where e.id=?", id)

	err := rows.Scan(&out.ID, &out.Name, &out.Dob, &out.City, &out.Majors, &out.Department.DeptID,
		&out.Department.DeptName, &out.Department.FloorNo)
	if err != nil {
		return entities.EmpDept{}, errors.EntityNotFound{Entity: "employee", ID: id.String()}
	}

	return out, nil
}

func (e Store) GetAll(ctx *gofr.Context, name string, include bool) (interface{}, error) {
	rows, err := e.GetRows(ctx, name)
	if err != nil {
		return []entities.EmpDept{}, err
	}

	data := make([]entities.EmpDept, 0)

	var temp entities.EmpDept

	for rows.Next() {
		err := rows.Scan(&temp.ID, &temp.Name, &temp.Dob, &temp.City, &temp.Majors, &temp.Department.DeptID,
			&temp.Department.DeptName, &temp.Department.FloorNo)
		if err != nil {
			return []entities.EmpDept{}, err
		}

		if name != "" && !include {
			temp.Department.FloorNo = 0
			temp.Department.DeptName = ""
		}

		data = append(data, temp)
	}

	if len(data) == 0 {
		return []entities.EmpDept{}, errors.DB{Err: fmt.Errorf("no data")}
	}

	return data, nil
}

func (e Store) GetRows(ctx *gofr.Context, name string) (*sql.Rows, error) {
	if name != "" {
		rows, err := ctx.DB().Query(
			"select e.id,e.name,e.dob,e.city,e.majors,d.id,d.name,d.floor from employee as e  "+
				"INNER JOIN department as d on e.DId=d.ID where e.name=?;", name)

		return rows, err
	}

	rows, err := ctx.DB().Query(
		"select e.id,e.name,e.dob,e.city,e.majors,d.id,d.name,d.floor from employee as e  INNER JOIN department as d on e.DId=d.id;")

	return rows, err
}

func (e Store) GetDepartment(ctx *gofr.Context, id int) (entities.Department, error) {
	var out entities.Department

	rows := ctx.DB().QueryRow("select * from department where id=?", id)

	err := rows.Scan(&out.DeptID, &out.DeptName, &out.FloorNo)
	if err != nil {
		return entities.Department{}, err
	}

	return out, nil
}
