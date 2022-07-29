package datastore

import (
	"Customer/models"
	"context"
	"database/sql"
)

type Store struct {
	db *sql.DB
}

func new(db *sql.DB) Store {
	return Store{db: db}
}

func (s Store) post(ctx context.Context, customer *models.Customer) (models.Customer, error) {
	_, err := s.db.ExecContext(ctx, "INSERT INTO Customers (firstName, lastName, dob, city) VALUES (?,?,?,?)", customer.FirstName, customer.LastName, customer.Dob, customer.City)
	if err != nil {
		return models.Customer{}, err
	}

	return *customer, nil
}
