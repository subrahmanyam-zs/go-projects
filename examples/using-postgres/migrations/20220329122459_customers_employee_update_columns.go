package migrations

import (
	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/log"
)

type K20220329122459 struct {
}

func (k K20220329122459) Up(d *datastore.DataStore, logger log.Logger) error {
	_, err := d.DB().Exec(AddCountry)
	if err != nil {
		return err
	}

	return nil
}

func (k K20220329122459) Down(d *datastore.DataStore, logger log.Logger) error {
	_, err := d.DB().Exec(DropCountry)
	if err != nil {
		return err
	}

	return nil
}
