package types

import (
	"regexp"

	"github.com/zopsmart/gofr/pkg/errors"
)

// this will compile the regex once instead of compiling it each time when it is being called.
var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~[:^ascii:]-]+@(?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9]\\.)+[a-zA-Z]{2,}$")

type Email string

func (e Email) Check() error {
	if !emailRegex.MatchString(string(e)) {
		return errors.InvalidParam{Param: []string{"emailAddress"}}
	}

	return nil
}
