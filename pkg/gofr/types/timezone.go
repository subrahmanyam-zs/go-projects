package types

import (
	"time"

	"developer.zopsmart.com/go/gofr/pkg/errors"
)

type TimeZone string

func (t TimeZone) Check() error {
	_, err := time.LoadLocation(string(t))
	if err != nil {
		return errors.InvalidParam{Param: []string{"timeZone"}}
	}

	return nil
}
