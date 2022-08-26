package errorsHandler

type NoData struct {
	Msg string
}

func (e NoData) Error() string {
	return e.Msg
}
