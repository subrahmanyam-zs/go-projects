package store

import (
	"database/sql"
	"fmt"

	"developer.zopsmart.com/go/gofr/examples/using-postgres/model"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Model struct{}

func New() Store {
	return &Model{}
}

type Store interface {
	Get(c *gofr.Context) (*[]model.Customer, error)
	GetByID(c *gofr.Context, id int) (*model.Customer, error)
	Update(c *gofr.Context, customer model.Customer) (*model.Customer, error)
	Create(c *gofr.Context, customer model.Customer) (*model.Customer, error)
	Delete(c *gofr.Context, id int) error
}

func (m Model) Get(c *gofr.Context) (*[]model.Customer, error) {
	rows, err := c.DB().QueryContext(c, "SELECT * FROM customers")
	if err != nil {
		return nil, errors.DB{Err: err}
	}

	defer func() {
		_ = rows.Close()
		_ = rows.Err() // or modify return value
	}()

	customers := make([]model.Customer, 0)

	for rows.Next() {
		var customer model.Customer

		err := rows.Scan(&customer.ID, &customer.Name)
		if err != nil {
			return nil, err
		}

		customers = append(customers, customer)
	}

	return &customers, nil
}

func (m Model) GetByID(c *gofr.Context, id int) (*model.Customer, error) {
	var customer model.Customer

	err := c.DB().QueryRowContext(c, " SELECT * FROM customers where id=$1", id).Scan(&customer.ID, &customer.Name)
	if err == sql.ErrNoRows {
		return nil, errors.EntityNotFound{
			Entity: "customer",
			ID:     fmt.Sprint(id),
		}
	}

	return &customer, nil
}

func (m Model) Update(c *gofr.Context, customer model.Customer) (*model.Customer, error) {
	_, err := c.DB().ExecContext(c, "UPDATE customers SET name=$1 WHERE id=$2", customer.Name, customer.ID)
	if err != nil {
		return nil, errors.DB{Err: err}
	}

	return &customer, nil
}

func (m Model) Create(c *gofr.Context, customer model.Customer) (*model.Customer, error) {
	var resp model.Customer

	err := c.DB().QueryRowContext(c, "INSERT INTO customers(name) VALUES($1) RETURNING id, name", customer.Name).Scan(
		&resp.ID, &resp.Name,
	)
	if err != nil {
		return nil, errors.DB{Err: err}
	}

	return &resp, nil
}

func (m Model) Delete(c *gofr.Context, id int) error {
	_, err := c.DB().ExecContext(c, "DELETE FROM customers where id=$1", id)
	if err != nil {
		return errors.DB{Err: err}
	}

	return nil
}
