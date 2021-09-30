package stores

import (
	"developer.zopsmart.com/go/gofr/examples/using-cassandra/models"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Person interface {
	Get(ctx *gofr.Context, filter models.Person) []models.Person
	Create(ctx *gofr.Context, data models.Person) ([]models.Person, error)
	Delete(ctx *gofr.Context, id string) error
	Update(ctx *gofr.Context, data models.Person) ([]models.Person, error)
}
