package types

import (
	"unicode"

	"github.com/zopsmart/gofr/pkg/errors"
)

type Enum struct {
	ValidValues []string `json:"validValues"`
	Value       string   `json:"value"`
	Parameter   string
}

func (e Enum) Check() error {
	// Enum values MUST be in the UPPER_SNAKE format
	for _, v := range e.Value {
		// allow underscores and numbers as enum values
		if v == 95 || unicode.IsDigit(v) {
			continue
		} else if !unicode.IsLetter(v) || !unicode.IsUpper(v) {
			return errors.InvalidParam{Param: []string{e.Parameter}}
		}
	}

	for _, v := range e.ValidValues {
		if v == e.Value {
			return nil
		}
	}

	return errors.InvalidParam{Param: []string{e.Parameter}}
}
