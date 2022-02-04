package store

import (
	"database/sql"
	"fmt"

	"developer.zopsmart.com/go/gofr/examples/using-postgres/model"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type customer struct{}

// New is factory function for store layer
func New() Store {
	return customer{}
}

type Store interface {
	Get(ctx *gofr.Context) ([]model.Customer, error)
	GetByID(ctx *gofr.Context, id int) (model.Customer, error)
	Update(ctx *gofr.Context, customer model.Customer) (model.Customer, error)
	Create(ctx *gofr.Context, customer model.Customer) (model.Customer, error)
	Delete(ctx *gofr.Context, id int) error
}

func (c customer) Get(ctx *gofr.Context) ([]model.Customer, error) {
	rows, err := ctx.DB().QueryContext(ctx, "SELECT * FROM customers")
	if err != nil {
		return nil, errors.DB{Err: err}
	}

	defer rows.Close()

	customers := make([]model.Customer, 0)

	for rows.Next() {
		var c model.Customer

		err = rows.Scan(&c.ID, &c.Name)
		if err != nil {
			return nil, errors.DB{Err: err}
		}

		customers = append(customers, c)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return customers, nil
}

func (c customer) GetByID(ctx *gofr.Context, id int) (model.Customer, error) {
	var resp model.Customer

	err := ctx.DB().QueryRowContext(ctx, " SELECT * FROM customers where id=$1", id).Scan(&resp.ID, &resp.Name)
	if err == sql.ErrNoRows {
		return model.Customer{}, errors.EntityNotFound{Entity: "customer", ID: fmt.Sprint(id)}
	}

	return resp, nil
}

func (c customer) Update(ctx *gofr.Context, cust model.Customer) (model.Customer, error) {
	_, err := ctx.DB().ExecContext(ctx, "UPDATE customers SET name=$1 WHERE id=$2", cust.Name, cust.ID)
	if err != nil {
		return model.Customer{}, errors.DB{Err: err}
	}

	return cust, nil
}

func (c customer) Create(ctx *gofr.Context, cust model.Customer) (model.Customer, error) {
	var resp model.Customer

	err := ctx.DB().QueryRowContext(ctx, "INSERT INTO customers(name) VALUES($1) RETURNING id, name", cust.Name).Scan(
		&resp.ID, &resp.Name,
	)
	if err != nil {
		return model.Customer{}, errors.DB{Err: err}
	}

	return resp, nil
}

func (c customer) Delete(ctx *gofr.Context, id int) error {
	_, err := ctx.DB().ExecContext(ctx, "DELETE FROM customers where id=$1", id)
	if err != nil {
		return errors.DB{Err: err}
	}

	return nil
}
