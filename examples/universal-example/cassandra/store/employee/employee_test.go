package employee

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/examples/universal-example/cassandra/entity"
	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/config"
	"developer.zopsmart.com/go/gofr/pkg/log"
)

func initializeTest(t *testing.T) *gofr.Gofr {
	c := config.NewGoDotEnvProvider(log.NewLogger(), "../../../configs")
	k := gofr.NewWithConfig(c)

	query := "CREATE TABLE IF NOT exists employees (id int, name text, phone text, email text, city text, PRIMARY KEY (id) )"

	err := k.Cassandra.Session.Query(query).Exec()
	if err != nil {
		k.Logger.Error("Employee table does not exist: ", err)
	}

	// initializing the seeder
	seeder := datastore.NewSeeder(&k.DataStore, "../../db")
	seeder.RefreshCassandra(t, "employees")

	return k
}

func TestCassandraEmployee_Get(t *testing.T) {
	tests := []struct {
		input  entity.Employee
		output []entity.Employee
	}{
		{entity.Employee{ID: 1}, []entity.Employee{{ID: 1, Name: "Rohan", Phone: "01222", Email: "rohan@zopsmart.com", City: "Berlin"}}},
		{entity.Employee{ID: 1, Name: "Rohan", Phone: "01222", Email: "rohan@zopsmart.com", City: "Berlin"},
			[]entity.Employee{{ID: 1, Name: "Rohan", Phone: "01222", Email: "rohan@zopsmart.com", City: "Berlin"}}},
		{entity.Employee{}, []entity.Employee{{ID: 1, Name: "Rohan", Phone: "01222", Email: "rohan@zopsmart.com", City: "Berlin"},
			{ID: 2, Name: "Aman", Phone: "22234", Email: "aman@zopsmart.com", City: "Florida"}}},
		{entity.Employee{ID: 7, Name: "Sunita"}, nil},
	}

	k := initializeTest(t)
	ctx := gofr.NewContext(nil, nil, k)

	for i, tc := range tests {
		output := New().Get(ctx, tc.input)
		assert.Equal(t, tc.output, output, i)
	}
}

func TestCassandraEmployee_Create(t *testing.T) {
	tests := []struct {
		input  entity.Employee
		output []entity.Employee
		err    error
	}{
		{entity.Employee{ID: 3, Name: "Sunita", Phone: "01234", Email: "sunita@zopsmart.com", City: "Darbhanga"},
			[]entity.Employee{{ID: 3, Name: "Sunita", Phone: "01234", Email: "sunita@zopsmart.com", City: "Darbhanga"}}, nil},
		{entity.Employee{ID: 4, Name: "Anna", Phone: "01333", Email: "anna@zopsmart.com", City: "Delhi"},
			[]entity.Employee{{ID: 4, Name: "Anna", Phone: "01333", Email: "anna@zopsmart.com", City: "Delhi"}}, nil},
	}

	k := initializeTest(t)
	ctx := gofr.NewContext(nil, nil, k)

	for i, tc := range tests {
		output, err := New().Create(ctx, tc.input)
		assert.Equal(t, tc.err, err, i)
		assert.Equal(t, tc.output, output, i)
	}
}
