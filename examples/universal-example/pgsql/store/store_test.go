package store

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/examples/universal-example/pgsql/entity"
	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func initializeTest(t *testing.T) *gofr.Gofr {
	app := gofr.New()
	queryTable := `
 	   CREATE TABLE IF NOT EXISTS employees (
	   id         int primary key,
	   name       varchar (50),
 	   phone      varchar(50),
 	   email      varchar(50) ,
 	   city       varchar(50))
	`

	if _, err := app.DB().Exec(queryTable); err != nil {
		app.Logger.Errorf("Postgres Got error while sourcing the schema: ", err)
	}
	// initializing the seeder
	seeder := datastore.NewSeeder(&app.DataStore, "../db")
	seeder.RefreshTables(t, "employees")

	return app
}

func TestEmployee_Get(t *testing.T) {
	app := initializeTest(t)
	c := gofr.NewContext(nil, nil, app)

	_, err := New().Get(c)
	assert.Equal(t, nil, err)
}

func TestEmployee_Create(t *testing.T) {
	app := initializeTest(t)
	c := gofr.NewContext(nil, nil, app)
	e := New()

	tests := []struct {
		desc    string
		input   entity.Employee
		wantErr bool
	}{
		{"success", entity.Employee{ID: 9, Name: "Sunita", Phone: "01234", Email: "sunita@zopsmart.com", City: "Darbhanga"}, false},
		{"fail-primary key constraint violation", entity.Employee{ID: 9, Name: "Angi", Phone: "01333", Email: "anna@zopsmart.com", City: "Delhi"},
			true},
	}

	for i, tc := range tests {
		err := e.Create(c, tc.input)
		if (err != nil) != tc.wantErr {
			t.Errorf("Testcase[%v] Failed in create", i)
		}
	}
}
