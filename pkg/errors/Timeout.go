package errors

import (
	"fmt"
)

type Timeout struct {
	URL string
}

func (t Timeout) Error() string {
	return fmt.Sprintf("Request to %v has Timed out!", t.URL)
}
