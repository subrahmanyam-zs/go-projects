package customer

import (
	"testing"

	"developer.zopsmart.com/go/gofr/examples/data-layer-with-mongo/models"
	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"github.com/stretchr/testify/assert"
)

func initializeTest(t *testing.T) *gofr.Gofr {
	app := gofr.New()

	// initializing the seeder
	seeder := datastore.NewSeeder(&app.DataStore, "../../db")
	seeder.RefreshMongoCollections(t, "customers")

	return app
}

func TestCustomer_Get(t *testing.T) {
	tests := []struct {
		desc string
		name string
		resp []models.Customer
		err  error
	}{
		{"get single entity", "Messi", []models.Customer{{Name: "Messi", Age: 32, City: "Barcelona"}}, nil},
		{"get multiple entity", "Tim", []models.Customer{{Name: "Tim", Age: 53, City: "London"}, {Name: "Tim", Age: 35, City: "Munich"}}, nil},
	}

	store := New()
	app := initializeTest(t)

	for i, tc := range tests {
		ctx := gofr.NewContext(nil, nil, app)
		resp, err := store.Get(ctx, tc.name)

		assert.Equal(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)

		assert.Equal(t, tc.resp, resp, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

func TestModel_Create(t *testing.T) {
	tests := []struct {
		desc     string
		customer string
		err      error
	}{
		{"create succuss", `{"name":"Pirlo","age":42,"city":"Turin"}`, nil},
	}

	store := New()
	app := initializeTest(t)

	for i, tc := range tests {
		ctx := gofr.NewContext(nil, nil, app)

		var c models.Customer

		err := store.Create(ctx, c)

		assert.Equal(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

func TestModel_Delete(t *testing.T) {
	tests := []struct {
		desc  string
		name  string
		count int
		err   error
	}{
		{"delete non existent entity", "Alex", 0, nil},
		{"delete multiple entity", "Tim", 2, nil},
		{"delete single entity", "Thomas", 1, nil},
	}

	store := New()
	app := initializeTest(t)

	for i, tc := range tests {
		ctx := gofr.NewContext(nil, nil, app)
		count, err := store.Delete(ctx, tc.name)

		assert.Equal(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)

		assert.Equal(t, tc.count, count, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}
