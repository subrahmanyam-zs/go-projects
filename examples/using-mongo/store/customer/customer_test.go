package customer

import (
	"reflect"
	"testing"

	"developer.zopsmart.com/go/gofr/examples/using-mongo/entity"
	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func initializeTest(t *testing.T) *gofr.Gofr {
	k := gofr.New()

	//initializing the seeder
	seeder := datastore.NewSeeder(&k.DataStore, "../../db")
	seeder.RefreshMongoCollections(t, "customers")

	return k
}

func TestCustomer_Get(t *testing.T) {
	//nolint: govet, table tests
	testCases := []struct {
		name         string
		expectedResp []*entity.Customer
	}{
		//{"Alex", nil},
		{"Messi", []*entity.Customer{{"Messi", 32, "Barcelona"}}},
		{"Tim", []*entity.Customer{{"Tim", 53, "London"}, {"Tim", 35, "Munich"}}},
	}

	customer := Customer{}
	k := initializeTest(t)

	for index, tc := range testCases {
		context2 := gofr.NewContext(nil, nil, k)
		resp, _ := customer.Get(context2, tc.name)

		if !reflect.DeepEqual(resp, tc.expectedResp) {
			t.Errorf("Testcase[%v] Failed\tExpected %v\nGot %v\n", index, tc.expectedResp, resp)
		}
	}
}

func TestModel_Create(t *testing.T) {
	testCases := []struct {
		customer string
		err      error
	}{
		{`{"name":"Pirlo","age":42,"city":"Turin"}`, nil},
	}

	customer := Customer{}
	k := initializeTest(t)

	for i, tc := range testCases {
		context2 := gofr.NewContext(nil, nil, k)

		var model entity.Customer

		err := customer.Create(context2, &model)

		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("Testcase[%v] Failed\tExpected %vGot %v\n", i, tc.err, err)
		}
	}
}

func TestModel_Delete(t *testing.T) {
	testCases := []struct {
		name                string
		expectedDeleteCount int
	}{
		{"Alex", 0},
		{"Tim", 2},
		{"Thomas", 1},
	}

	customer := Customer{}
	k := initializeTest(t)

	for i, tc := range testCases {
		context2 := gofr.NewContext(nil, nil, k)
		deleteCount, _ := customer.Delete(context2, tc.name)

		if deleteCount != tc.expectedDeleteCount {
			t.Errorf("Testcase[%v] failed\nExpected delete count: %v \nGot delete count: %v",
				i,
				tc.expectedDeleteCount,
				deleteCount,
			)
		}
	}
}
