package types

import (
	"regexp"

	"developer.zopsmart.com/go/gofr/pkg/errors"
)

// nolint:lll // this will compile the regex once instead of compiling it each time when it is being called.
var durationRegex = regexp.MustCompile(`^P((((\d+(\.\d+)?)Y)?((\d+(\.\d+)?)M)?((\d+(\.\d+)?)D)?)?(T((\d+(\.\d+)?)H)?((\d+(\.\d+)?)M)?((\d+(\.\d+)?)S)?)?){1}$|(^P(\d+(\.\d+)?)W$)`)

type Duration string

func (d Duration) Check() error {
	// The format for duration is :PnYnMnDTnHnMnS or PnW
	const durationLen = 3
	if len(d) < durationLen {
		return errors.InvalidParam{Param: []string{"duration"}}
	}

	matches := durationRegex.FindStringSubmatch(string(d))

	if len(matches) == 0 {
		return errors.InvalidParam{Param: []string{"duration"}}
	}

	return nil
}
