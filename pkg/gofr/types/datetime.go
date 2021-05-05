package types

import (
	"time"

	"github.com/zopsmart/gofr/pkg/errors"
)

type Datetime struct {
	Value    string `json:"value"`
	Timezone string `json:"timezone"`
}

func (d Datetime) Check() error {
	// date and time together MUST always be included in the datetime structure
	_, err := time.Parse(time.RFC3339, d.Value)
	if err != nil {
		return errors.InvalidParam{Param: []string{"datetime"}}
	}

	err = Validate(TimeZone(d.Timezone))
	if err != nil {
		return errors.InvalidParam{Param: []string{"timezone"}}
	}

	return nil
}
