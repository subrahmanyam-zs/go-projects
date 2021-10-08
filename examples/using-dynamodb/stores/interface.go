package stores

import (
	"developer.zopsmart.com/go/gofr/examples/using-dynamodb/models"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Person interface {
	Get(c *gofr.Context, id string) (models.Person, error)
	Create(c *gofr.Context, user models.Person) error
	Update(c *gofr.Context, user models.Person) error
	Delete(c *gofr.Context, id string) error
}
