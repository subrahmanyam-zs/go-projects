package brand

import (
	"developer.zopsmart.com/go/gofr/examples/mock-c-layer/models"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type brand struct{}

// New is factory function for store layer
//nolint:revive // brand should not be used without proper initilization with required dependency
func New() brand {
	return brand{}
}

func (b brand) Get(ctx *gofr.Context) ([]models.Brand, error) {
	id := ctx.Param("id")

	const (
		id1 = 1
		id2 = 2
	)

	switch id {
	case "1":
		return []models.Brand{{ID: id1, Name: "brand 1"}}, nil
	case "2":
		return []models.Brand{{ID: id1, Name: "brand 1"}, {ID: id2, Name: "brand 2"}}, nil
	}

	return nil, errors.Error("core error")
}

func (b brand) Create(ctx *gofr.Context, brand models.Brand) (models.Brand, error) {
	if brand.Name == "brand 3" {
		return models.Brand{}, errors.Error("core error")
	}

	return brand, nil
}
