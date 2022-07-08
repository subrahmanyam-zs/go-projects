package migrations

import (
	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/log"
)

type K20220329123813 struct {
}

func (k K20220329123813) Up(d *datastore.DataStore, logger log.Logger) error {
	_, err := d.DB().Exec(AddNotNullColumn)
	if err != nil {
		return err
	}

	return nil
}

func (k K20220329123813) Down(d *datastore.DataStore, logger log.Logger) error {
	_, err := d.DB().Exec(DeleteNotNullColumn)
	if err != nil {
		return err
	}

	return nil
}
