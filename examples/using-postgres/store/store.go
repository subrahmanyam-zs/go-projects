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
	customers := make([]model.Customer, 0)

	var customer model.Customer
	rows, err := c.DB().Query("SELECT * FROM customers")
	if err != nil {
		return nil, errors.DB{Err: err}
	}

	for rows.Next() {
		rows.Scan(&customer.ID, &customer.Name)
		customers = append(customers, customer)
	}

	return &customers, nil
}

func (m Model) GetByID(c *gofr.Context, id int) (*model.Customer, error) {
	var customer model.Customer

	err := c.DB().QueryRow(" SELECT * FROM customers where id=$1", id).Scan(&customer.ID, &customer.Name)

	if err == sql.ErrNoRows {
		return nil, errors.EntityNotFound{
			Entity: "customer",
			ID:     fmt.Sprint(id),
		}
	}

	return &customer, nil
}

func (m Model) Update(c *gofr.Context, customer model.Customer) (*model.Customer, error) {
	_, err := c.DB().Exec("UPDATE customers SET name = $1 WHERE id = $2", customer.Name, customer.ID)
	if err != nil {
		return nil, errors.DB{Err: err}
	}

	return m.GetByID(c, customer.ID)
}

func (m Model) Create(c *gofr.Context, customer model.Customer) (*model.Customer, error) {
	_, err := c.DB().Exec("INSERT INTO customers(name) VALUES($1)", customer.Name)

	if err != nil {
		return nil, errors.DB{Err: err}
	}

	var lastID int
	err = c.DB().QueryRow(`SELECT MAX(ID) AS MAX_ID FROM customers`).Scan(&lastID)
	return m.GetByID(c, lastID)
}

func (m Model) Delete(c *gofr.Context, id int) error {
	_, err := c.DB().Exec("DELETE FROM customers where id=$1", id)
	if err != nil {
		return errors.DB{Err: err}
	}

	return nil
}
