package types

import (
	"strconv"
	"strings"

	"developer.zopsmart.com/go/gofr/pkg/errors"

	"golang.org/x/text/currency"
)

type Currency string

func (c Currency) Check() error {
	// Currencies MUST use the ISO 4217 currency codes. Ex: USD 34.55
	currencyArray := strings.Fields(string(c))

	const currencyArrayLen = 2
	if len(currencyArray) != currencyArrayLen {
		return errors.InvalidParam{Param: []string{"currency"}}
	}

	_, err := currency.ParseISO(currencyArray[0])
	if err != nil {
		return errors.InvalidParam{Param: []string{"currencyCountryCode"}}
	}

	_, err = strconv.ParseFloat(currencyArray[1], 64)
	if err != nil {
		return errors.InvalidParam{Param: []string{"currencyValue"}}
	}

	return nil
}
