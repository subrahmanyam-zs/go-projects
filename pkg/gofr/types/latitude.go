package types

import (
	"github.com/zopsmart/gofr/pkg/errors"
)

type Latitude float64

func (l *Latitude) Check() error {
	if *l > 90 || *l < -90 {
		return errors.InvalidParam{Param: []string{"lat"}}
	}

	return nil
}
