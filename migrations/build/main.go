// This is auto-generated file using 'gofr migrate' tool. DO NOT EDIT.
package main

import (
	"developer.zopsmart.com/go/gofr/migrations"
	"developer.zopsmart.com/go/gofr/cmd/gofr/migration"
	"developer.zopsmart.com/go/gofr/cmd/gofr/migration/dbMigration"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	k := gofr.New()
	db := dbmigration.NewGorm(k.GORM())	

	err := migration.Migrate("cnsr-gofr", db, migrations.All(), "UP", k.Logger)
	if err != nil {
		k.Logger.Error(err)
	}
}
