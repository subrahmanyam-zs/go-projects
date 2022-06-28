package errors

import "fmt"

type DataStoreNotInitialized struct {
	DBName string
	Reason string
}

func (d DataStoreNotInitialized) Error() string {
	return fmt.Sprintf("couldn't initialize %v, %v", d.DBName, d.Reason)
}
