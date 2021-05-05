package store

import (
	"reflect"
	"testing"

	"github.com/zopsmart/gofr/examples/using-postgres/model"
	"github.com/zopsmart/gofr/pkg/datastore"
	"github.com/zopsmart/gofr/pkg/errors"
	"github.com/zopsmart/gofr/pkg/gofr"
)

func TestCoreLayer(t *testing.T) {
	k := gofr.New()

	// initializing the seeder
	seeder := datastore.NewSeeder(&k.DataStore, "../db")
	seeder.ResetCounter = true

	testAddCustomer(t, k, seeder)
	testGetCustomerByID(t, k, seeder)
	testUpdateCustomer(t, k, seeder)
	testGetCustomers(t, k, seeder)
	testDeleteCustomer(t, k, seeder)
}

func testAddCustomer(t *testing.T, k *gofr.Gofr, seeder *datastore.Seeder) {
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
		{
			customer: model.Customer{
				Name: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			},
			expectedErr: errors.DB{},
		},
	}

	seeder.RefreshTables(t, "customers")
	for i, test := range tests {
		c := gofr.NewContext(nil, nil, k)

		_, err := Model{}.Create(c, test.customer)

		d, _ := Model{}.Get(c)
		k.Logger.Log(d)

		if _, ok := err.(errors.DB); err != nil && ok == false {
			t.Errorf("Testcase %v FAILED", i)
		}
		if err != test.expectedErr && test.expectedErr == nil {
			t.Errorf("Testcase %v FAILED, got err: %v expected: %v", i, err, test.expectedErr)
		}
	}
}

func testGetCustomerByID(t *testing.T, k *gofr.Gofr, seeder *datastore.Seeder) {
	seeder.RefreshTables(t, "customers")
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

		_, err := Model{}.GetByID(c, test.id)
		if !reflect.DeepEqual(err, test.err) {
			t.Errorf("FAILED, Expected: %v, Got: %v", test.err, err)
		}
	}
}

func testUpdateCustomer(t *testing.T, k *gofr.Gofr, seeder *datastore.Seeder) {
	seeder.RefreshTables(t, "customers")
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
				Name: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			},
			expectedErr: errors.DB{},
		},
	}

	for i, test := range tests {
		c := gofr.NewContext(nil, nil, k)

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

func testGetCustomers(t *testing.T, k *gofr.Gofr, seeder *datastore.Seeder) {
	seeder.RefreshTables(t, "customers")
	c := gofr.NewContext(nil, nil, k)

	_, err := Model{}.Get(c)
	if err != nil {
		t.Errorf("FAILED, Expected: %v, Got: %v", nil, err)
	}
}

func testDeleteCustomer(t *testing.T, k *gofr.Gofr, seeder *datastore.Seeder) {
	seeder.RefreshTables(t, "customers")
	c := gofr.NewContext(nil, nil, k)

	err := Model{}.Delete(c, 12)
	if err != nil {
		t.Errorf("FAILED, Expected: %v, Got: %v", nil, err)
	}
}

func TestNew(t *testing.T) {
	_ = New()
}
