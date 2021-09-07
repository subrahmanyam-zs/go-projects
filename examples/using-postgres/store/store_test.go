package store

import (
	"context"
	"reflect"
	"testing"

	"developer.zopsmart.com/go/gofr/examples/using-postgres/model"
	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func TestCoreLayer(t *testing.T) {
	k := gofr.New()

	// initializing the seeder
	seeder := datastore.NewSeeder(&k.DataStore, "../db")
	seeder.ResetCounter = true

	createTable(k)
	seeder.RefreshTables(t, "customers")
	testGetCustomerByID(t, k)
	testAddCustomer(t, k)
	testAddCustomerWithError(t, k)
	testUpdateCustomer(t, k)
	testGetCustomers(t, k)
	testDeleteCustomer(t, k)
}

func createTable(k *gofr.Gofr) {
	_, err := k.DB().Exec("CREATE TABLE customers (id serial primary key,name varchar (50));")
	if err != nil {
		return
	}
}

func testAddCustomer(t *testing.T, k *gofr.Gofr) {
	tests := []struct {
		customer    model.Customer
		expectedErr error
	}{
		{
			customer: model.Customer{
				Name: "Test123",
			},
			expectedErr: nil,
		},
		{
			customer: model.Customer{
				Name: "Test234",
			},
			expectedErr: nil,
		},
	}

	for i, test := range tests {
		c := gofr.NewContext(nil, nil, k)
		c.Context = context.Background()

		resp, err := Model{}.Create(c, test.customer)

		k.Logger.Log(resp)

		if err != test.expectedErr && test.expectedErr == nil {
			t.Errorf("Testcase %v FAILED, got err: %v expected: %v", i, err, test.expectedErr)
		}
	}
}

func testAddCustomerWithError(t *testing.T, k *gofr.Gofr) {
	customer := model.Customer{
		Name: "very-long-mock-name-aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}

	c := gofr.NewContext(nil, nil, k)
	c.Context = context.Background()
	_, err := Model{}.Create(c, customer)

	if _, ok := err.(errors.DB); err != nil && ok == false {
		t.Errorf("Error Testcase FAILED")
	}
}

func testGetCustomerByID(t *testing.T, k *gofr.Gofr) {
	tests := []struct {
		id  int
		err error
	}{
		{
			id:  1,
			err: nil,
		},
		{
			id: 1223,
			err: errors.EntityNotFound{
				Entity: "customer",
				ID:     "1223",
			},
		},
	}

	for _, test := range tests {
		c := gofr.NewContext(nil, nil, k)
		c.Context = context.Background()

		_, err := Model{}.GetByID(c, test.id)
		if !reflect.DeepEqual(err, test.err) {
			t.Errorf("FAILED, Expected: %v, Got: %v", test.err, err)
		}
	}
}

func testUpdateCustomer(t *testing.T, k *gofr.Gofr) {
	tests := []struct {
		customer    model.Customer
		expectedErr error
	}{
		{
			customer: model.Customer{
				ID:   1,
				Name: "Test1234",
			},
			expectedErr: nil,
		},
		{
			customer: model.Customer{
				ID:   1,
				Name: "very-long-mock-name-aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			},
			expectedErr: errors.DB{},
		},
	}

	for i, test := range tests {
		c := gofr.NewContext(nil, nil, k)
		c.Context = context.Background()

		_, err := Model{}.Update(c, test.customer)
		if test.expectedErr == nil {
			if err != nil {
				t.Errorf("Testcase %v FAILED", i)
			}
		} else {
			if _, ok := err.(errors.DB); ok == false {
				t.Errorf("Testcase %v FAILED", i)
			}
		}
	}
}

func testGetCustomers(t *testing.T, k *gofr.Gofr) {
	c := gofr.NewContext(nil, nil, k)
	c.Context = context.Background()

	_, err := Model{}.Get(c)
	if err != nil {
		t.Errorf("FAILED, Expected: %v, Got: %v", nil, err)
	}
}

func testDeleteCustomer(t *testing.T, k *gofr.Gofr) {
	c := gofr.NewContext(nil, nil, k)
	c.Context = context.Background()

	err := Model{}.Delete(c, 12)
	if err != nil {
		t.Errorf("FAILED, Expected: %v, Got: %v", nil, err)
	}
}

func TestNew(t *testing.T) {
	_ = New()
}
