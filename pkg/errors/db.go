package errors

type DB struct {
	Err error
}

func (e DB) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}

	return "DB Error"
}
