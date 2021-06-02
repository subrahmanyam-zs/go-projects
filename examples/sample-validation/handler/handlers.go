package handler

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/types"
)

type Details struct {
	Phone types.Phone `json:"phone"`
	Email types.Email `json:"email"`
}

func ValidateEntry(c *gofr.Context) (interface{}, error) {
	var detail Details

	err := c.Bind(&detail)
	if err != nil {
		return nil, err
	}

	err = types.Validate(detail.Phone, detail.Email)
	if err != nil {
		return nil, err
	}

	return "successful validation", nil
}
