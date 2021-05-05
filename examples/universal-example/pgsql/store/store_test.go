package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zopsmart/gofr/examples/universal-example/pgsql/entity"
	"github.com/zopsmart/gofr/pkg/datastore"
	"github.com/zopsmart/gofr/pkg/gofr"
)

func initializeTest(t *testing.T) *gofr.Gofr {
	k := gofr.New()
	queryTable := `
 	   CREATE TABLE IF NOT EXISTS employees (
	   id         int primary key,
	   name       varchar (50),
 	   phone      varchar(50),
 	   email      varchar(50) ,
 	   city       varchar(50))
	`

	if _, err := k.DB().Exec(queryTable); err != nil {
		k.Logger.Errorf("Postgres Got error while sourcing the schema: ", err)
	}
	// initializing the seeder
	seeder := datastore.NewSeeder(&k.DataStore, "../db")
	seeder.RefreshTables(t, "employees")

	return k
}

func TestEmployee_Get(t *testing.T) {
	k := initializeTest(t)
	c := gofr.NewContext(nil, nil, k)

	_, err := New().Get(c)
	assert.Equal(t, nil, err)
}

func TestEmployee_Create(t *testing.T) {
	k := initializeTest(t)
	c := gofr.NewContext(nil, nil, k)
	{
		// Success Case
		emp := entity.Employee{ID: 9, Name: "Sunita", Phone: "01234", Email: "sunita@zopsmart.com", City: "Darbhanga"}
		err := New().Create(c, emp)
		assert.Equal(t, nil, err)
	}
	{
		// Failure Case using primary key constraint violation
		emp := entity.Employee{ID: 9, Name: "Angi", Phone: "01333", Email: "anna@zopsmart.com", City: "Delhi"}
		err := New().Create(c, emp)
		if err == nil {
			t.Error("Test Failed in create")
		}
	}
}
