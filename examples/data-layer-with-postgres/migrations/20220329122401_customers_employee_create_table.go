package migrations

import (
	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/log"
)

type K20220329122401 struct {
}

func (k K20220329122401) Up(d *datastore.DataStore, logger log.Logger) error {
	_, err := d.DB().Exec(CreateTable)
	if err != nil {
		return err
	}

	return nil
}

func (k K20220329122401) Down(d *datastore.DataStore, logger log.Logger) error {
	_, err := d.DB().Exec(DroopTable)
	if err != nil {
		return err
	}

	return nil
}
