package handlers

import (
	"fmt"

	"github.com/zopsmart/gofr/examples/using-mongo/entity"
	"github.com/zopsmart/gofr/examples/using-mongo/store"
	"github.com/zopsmart/gofr/pkg/gofr"
)

type Customer struct {
	model store.Customer
}

func New(c store.Customer) Customer {
	return Customer{
		model: c,
	}
}

func (cm Customer) Get(c *gofr.Context) (interface{}, error) {
	name := c.Param("name")

	results, err := cm.model.Get(c, name)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (cm Customer) Create(c *gofr.Context) (interface{}, error) {
	var model entity.Customer

	err := c.Bind(&model)
	if err != nil {
		return nil, err
	}

	err = cm.model.Create(c, &model)
	if err != nil {
		return nil, err
	}

	return "New Customer Added!", nil
}

func (cm Customer) Delete(c *gofr.Context) (interface{}, error) {
	name := c.Param("name")

	deleteCount, err := cm.model.Delete(c, name)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("%v Customers Deleted!", deleteCount), nil
}
