package migrations

import (
	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/log"
)

type K20220329122659 struct {
}

func (k K20220329122659) Up(d *datastore.DataStore, logger log.Logger) error {
	_, err := d.DB().Exec(AlterType)
	if err != nil {
		return err
	}

	return nil
}

func (k K20220329122659) Down(d *datastore.DataStore, logger log.Logger) error {
	_, err := d.DB().Exec(ResetType)
	if err != nil {
		return err
	}

	return nil
}
