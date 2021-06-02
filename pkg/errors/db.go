package errors

type DB struct {
	Err error
}

func (e DB) Error() string {
	return e.Err.Error()
}
