package store

import (
	"context"
	"testing"

	"developer.zopsmart.com/go/gofr/examples/data-layer-with-postgres/model"
	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"github.com/stretchr/testify/assert"
)

func TestCoreLayer(t *testing.T) {
	app := gofr.New()

	// initializing the seeder
	seeder := datastore.NewSeeder(&app.DataStore, "../db")
	seeder.ResetCounter = true

	createTable(app)
	seeder.RefreshTables(t, "customers")
	testGetCustomerByID(t, app)
	testAddCustomer(t, app)
	testAddCustomerWithError(t, app)
	testUpdateCustomer(t, app)
	testGetCustomers(t, app)
	testDeleteCustomer(t, app)
	testErrors(t, app)
}

func createTable(app *gofr.Gofr) {
	// drop table to clean previously added id's
	_, err := app.DB().Exec("DROP TABLE customers;")
	if err != nil {
		return
	}

	_, err = app.DB().Exec("CREATE TABLE customers (id serial primary key,name varchar (50));")
	if err != nil {
		return
	}
}

func testAddCustomer(t *testing.T, app *gofr.Gofr) {
	tests := []struct {
		desc     string
		customer model.Customer
		err      error
	}{
		{"create succuss test #1", model.Customer{Name: "Test123"}, nil},
		{"create succuss test #2", model.Customer{Name: "Test234"}, nil},
	}

	for i, tc := range tests {
		ctx := gofr.NewContext(nil, nil, app)
		ctx.Context = context.Background()

		store := New()
		resp, err := store.Create(ctx, tc.customer)

		app.Logger.Log(resp)

		assert.Equal(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

func testAddCustomerWithError(t *testing.T, app *gofr.Gofr) {
	customer := model.Customer{
		Name: "very-long-mock-name-lasdjflsdjfljasdlfjsdlfjsdfljlkj",
	}

	ctx := gofr.NewContext(nil, nil, app)
	ctx.Context = context.Background()

	store := New()

	_, err := store.Create(ctx, customer)
	if _, ok := err.(errors.DB); err != nil && ok == false {
		t.Errorf("Error Testcase FAILED")
	}
}

func testGetCustomerByID(t *testing.T, app *gofr.Gofr) {
	tests := []struct {
		desc string
		id   int
		err  error
	}{
		{"Get existent id", 1, nil},
		{"Get non existent id", 1223, errors.EntityNotFound{Entity: "customer", ID: "1223"}},
	}

	for i, tc := range tests {
		ctx := gofr.NewContext(nil, nil, app)
		ctx.Context = context.Background()

		store := New()

		_, err := store.GetByID(ctx, tc.id)
		assert.Equal(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

func testUpdateCustomer(t *testing.T, app *gofr.Gofr) {
	tests := []struct {
		desc     string
		customer model.Customer
		err      error
	}{
		{"update succuss", model.Customer{ID: 1, Name: "Test1234"}, nil},
		{"update fail", model.Customer{ID: 1, Name: "very-long-mock-name-lasdjflsdjfljasdlfjsdlfjsdfljlkj"}, errors.DB{}},
	}

	for i, tc := range tests {
		ctx := gofr.NewContext(nil, nil, app)
		ctx.Context = context.Background()

		store := New()

		_, err := store.Update(ctx, tc.customer)
		if _, ok := err.(errors.DB); err != nil && ok == false {
			t.Errorf("TEST[%v] Failed.\tExpected %v\tGot %v\n%s", i, tc.err, err, tc.desc)
		}
	}
}

func testGetCustomers(t *testing.T, app *gofr.Gofr) {
	ctx := gofr.NewContext(nil, nil, app)
	ctx.Context = context.Background()

	store := New()

	_, err := store.Get(ctx)
	if err != nil {
		t.Errorf("FAILED, Expected: %v, Got: %v", nil, err)
	}
}

func testDeleteCustomer(t *testing.T, app *gofr.Gofr) {
	ctx := gofr.NewContext(nil, nil, app)
	ctx.Context = context.Background()

	store := New()

	err := store.Delete(ctx, 12)
	if err != nil {
		t.Errorf("FAILED, Expected: %v, Got: %v", nil, err)
	}
}

func testErrors(t *testing.T, app *gofr.Gofr) {
	ctx := gofr.NewContext(nil, nil, app)
	ctx.Context = context.Background()
	_ = ctx.DB().Close() // close the connection to generate errors

	store := New()

	err := store.Delete(ctx, 12)
	if err == nil {
		t.Errorf("FAILED, Expected: %v, Got: %v", nil, err)
	}

	_, err = store.Get(ctx)
	if err == nil {
		t.Errorf("FAILED, Expected: %v, Got: %v", nil, err)
	}
}

func TestNew(t *testing.T) {
	_ = New()
}
