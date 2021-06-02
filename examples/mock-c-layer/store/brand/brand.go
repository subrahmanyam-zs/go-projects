package brand

import (
	"errors"

	"developer.zopsmart.com/go/gofr/examples/mock-c-layer/store"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Brand struct{}

func New() *Brand {
	return &Brand{}
}

func (b *Brand) Get(ctx *gofr.Context) ([]store.Model, error) {
	id := ctx.Param("id")

	switch id {
	case "1":
		return []store.Model{{ID: 1, Name: "brand 1"}}, nil
	case "2":
		return []store.Model{{ID: 1, Name: "brand 1"}, {ID: 2, Name: "brand 2"}}, nil
	}

	return nil, errors.New("core error")
}

func (b *Brand) Create(ctx *gofr.Context, brand store.Model) (store.Model, error) {
	if brand.Name == "brand 3" {
		return store.Model{}, errors.New("core error")
	}

	return brand, nil
}
