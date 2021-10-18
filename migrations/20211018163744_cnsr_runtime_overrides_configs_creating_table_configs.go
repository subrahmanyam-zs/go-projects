package migrations

import (
	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/log"
)

type K20211018163744 struct {
}

func (k K20211018163744) Up(d *datastore.DataStore, logger log.Logger) error {
	var val int

	err := d.DB().QueryRow("SELECT 2+2").Scan(&val)
	if err != nil {
		return err
	}

	return nil
}

func (k K20211018163744) Down(d *datastore.DataStore, logger log.Logger) error {
	return nil
}
