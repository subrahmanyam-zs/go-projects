package employee

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/google/uuid"

	entities "EmployeeDepartment/entities"
	"EmployeeDepartment/errorsHandler"
	"EmployeeDepartment/store"
)

type Store struct {
	Db *sql.DB
}

func New(db *sql.DB) Store {
	return Store{Db: db}
}

func (s Store) Create(ctx context.Context, emp *entities.Employee) (*entities.Employee, error) {
	var uid uuid.UUID = uuid.New()

	res, err := s.Db.ExecContext(ctx, "Insert into employee values(?,?,?,?,?,?)", uid, emp.Name, emp.Dob, emp.City, emp.Majors, emp.DId)
	if err != nil {
		return &entities.Employee{}, err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 1 {
		emp.ID = uid
		return emp, nil
	}

	return &entities.Employee{}, errorsHandler.AlreadyExists{Msg: "Already Exists"}
}

func (s Store) Update(ctx context.Context, id uuid.UUID, emp *entities.Employee) (*entities.Employee, error) {
	res, err := s.Db.ExecContext(ctx, "update employee set id=?, name=?,dob=?,city=?,majors=?,dId=? "+
		"where id=?;", emp.ID, emp.Name, emp.Dob, emp.City, emp.Majors, emp.DId, id)
	if err != nil {
		return &entities.Employee{}, err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 1 {
		emp.ID = id
		return emp, nil
	}

	return &entities.Employee{}, &errorsHandler.AlreadyExists{Msg: "Already Exists"}
}

func (s Store) Delete(ctx context.Context, id uuid.UUID) (int, error) {
	res, err := s.Db.ExecContext(ctx, "Delete from employee where id=?", id)
	if err != nil {
		return http.StatusBadRequest, err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 1 {
		return http.StatusNoContent, nil
	}

	return http.StatusBadRequest, &errorsHandler.IDNotFound{Msg: "ID not found"}
}

func (s Store) Read(ctx context.Context, id uuid.UUID) (entities.EmployeeAndDepartment, error) {
	var out entities.EmployeeAndDepartment

	rows := s.Db.QueryRowContext(ctx,
		"select e.id,e.name,e.dob,e.city,e.majors,d.id,d.name,d.floor from employee"+
			" as e inner join department as d on e.DId=d.id where e.id=?", id)

	err := rows.Scan(&out.ID, &out.Name, &out.Dob, &out.City, &out.Majors, &out.Dept.ID,
		&out.Dept.Name, &out.Dept.FloorNo)
	if err != nil {
		return entities.EmployeeAndDepartment{}, err
	}

	return out, nil
}

func (s Store) ReadAll(para store.Parameters) ([]entities.EmployeeAndDepartment, error) {
	rows, err := s.GetRows(para)
	if err != nil {
		return []entities.EmployeeAndDepartment{}, err
	}

	data := make([]entities.EmployeeAndDepartment, 0)

	var temp entities.EmployeeAndDepartment

	for rows.Next() {
		err = rows.Scan(&temp.ID, &temp.Name, &temp.Dob, &temp.City, &temp.Majors, &temp.Dept.ID, &temp.Dept.Name, &temp.Dept.FloorNo)
		if err != nil {
			return []entities.EmployeeAndDepartment{}, err
		}

		if para.Name != "" && !para.IncludeDepartment {
			temp.Dept.FloorNo = 0
			temp.Dept.Name = ""
		}

		data = append(data, temp)
	}

	if len(data) == 0 {
		return nil, errorsHandler.NoData{Msg: "No Data"}
	}

	return data, nil
}

func (s Store) GetRows(para store.Parameters) (*sql.Rows, error) {
	if para.Name != "" {
		rows, err := s.Db.QueryContext(para.Ctx,
			"select e.id,e.name,e.dob,e.city,e.majors,d.id,d.name,d.floor from employee as e  "+
				"INNER JOIN department as d on e.DId=d.ID where e.name=?;", para.Name)

		return rows, err
	}

	rows, err := s.Db.Query(
		"select e.id,e.name,e.dob,e.city,e.majors,d.id,d.name,d.floor from employee as e  INNER JOIN department as d on e.DId=d.id;")

	return rows, err
}

func (s Store) ReadDepartment(ctx context.Context, id int) (entities.Department, error) {
	var out entities.Department

	rows, err := s.Db.Query("select * from department where id=?", id)

	if err != nil {
		return entities.Department{}, err
	}

	rows.Next()

	err = rows.Scan(&out.ID, &out.Name, &out.FloorNo)
	if err != nil {
		return entities.Department{}, err
	}

	rows.Close()

	return out, nil
}
