package types

import (
	"time"

	"developer.zopsmart.com/go/gofr/pkg/errors"
)

type Time string

func (t Time) Check() error {
	_, err := time.Parse("15:04:05.000", string(t))
	if err != nil {
		return errors.InvalidParam{Param: []string{"time"}}
	}

	return nil
}
