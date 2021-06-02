package types

import (
	"time"

	"developer.zopsmart.com/go/gofr/pkg/errors"
)

type Date string

func (d Date) Check() error {
	_, err := time.Parse("2006-01-02", string(d))
	if err != nil {
		return errors.InvalidParam{Param: []string{"date"}}
	}

	return nil
}
