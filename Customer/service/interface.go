package service

import (
	"Customer/models"
	"context"
)

type Customer interface {
	Post(context.Context, *models.Customer) (models.Customer, error)
}
