package migrations

import (
	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/log"
)

type K20220329122607 struct {
}

func (k K20220329122607) Up(d *datastore.DataStore, logger log.Logger) error {
	_, err := d.DB().Exec(DropPhone)
	if err != nil {
		return err
	}

	return nil
}

func (k K20220329122607) Down(d *datastore.DataStore, logger log.Logger) error {
	_, err := d.DB().Exec(AddPhone)
	if err != nil {
		return err
	}

	return nil
}
