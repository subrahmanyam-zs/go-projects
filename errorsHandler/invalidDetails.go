package errorsHandler

import "fmt"

type InvalidDetails struct {
	Msg string
}

func (e InvalidDetails) Error() string {
	return fmt.Sprintf("Invalid %s", e.Msg)
}
