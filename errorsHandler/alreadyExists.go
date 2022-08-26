package errorsHandler

type AlreadyExists struct {
	Msg string
}

func (e AlreadyExists) Error() string {
	return e.Msg
}
