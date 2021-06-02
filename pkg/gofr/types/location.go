package types

import (
	"developer.zopsmart.com/go/gofr/pkg/errors"
)

type Location struct {
	Latitude  *Latitude  `json:"lat"`
	Longitude *Longitude `json:"lng"`
}

func (l Location) Check() error {
	if l.Latitude == nil && l.Longitude == nil {
		return errors.MultipleErrors{Errors: []error{errors.InvalidParam{Param: []string{"lat is nil"}},
			errors.InvalidParam{Param: []string{"lng is nil"}}}}
	}

	if l.Latitude == nil {
		return errors.InvalidParam{Param: []string{"lat is nil"}}
	}

	err := Validate(l.Latitude)
	if err != nil {
		return errors.InvalidParam{Param: []string{"lat"}}
	}

	if l.Longitude == nil {
		return errors.InvalidParam{Param: []string{"lng is nil"}}
	}

	err = Validate(l.Longitude)
	if err != nil {
		return errors.InvalidParam{Param: []string{"lng"}}
	}

	return nil
}
