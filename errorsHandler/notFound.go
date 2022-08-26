package errorsHandler

type IDNotFound struct {
	Msg string
}

func (e IDNotFound) Error() string {
	return e.Msg
}
