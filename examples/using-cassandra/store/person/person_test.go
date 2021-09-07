package person

import (
	"os"
	"reflect"
	"testing"

	"developer.zopsmart.com/go/gofr/examples/using-cassandra/entity"
	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func TestMain(m *testing.M) {
	k := gofr.New()
	// Create a table person if the table does not exists
	queryStr := "CREATE TABLE IF NOT EXISTS persons (id int PRIMARY KEY, name text, age int, state text )"
	err := k.Cassandra.Session.Query(queryStr).Exec()
	// if table creation is unsuccessful log the error
	if err != nil {
		k.Logger.Errorf("Failed creation of table persons :%v", err)
	} else {
		k.Logger.Info("Table persons created Successfully")
	}

	os.Exit(m.Run())
}

func initializeTest(t *testing.T) *gofr.Gofr {
	k := gofr.New()
	// initializing the seeder
	seeder := datastore.NewSeeder(&k.DataStore, "../../db")
	seeder.RefreshCassandra(t, "persons")

	return k
}

func createMap(input []*entity.Person) map[entity.Person]int {
	output := make(map[entity.Person]int)

	for _, val := range input {
		if val != nil {
			output[*val]++
		}
	}

	return output
}

func isSubset(supSet, subSet []*entity.Person) bool {
	set := createMap(supSet)
	subset := createMap(subSet)

	for k := range subset {
		if val, ok := set[k]; !ok || val != subset[k] {
			return false
		}
	}

	return true
}

func TestGet(t *testing.T) {
	tests := []struct {
		input  entity.Person
		output []*entity.Person
		err    error
	}{
		{entity.Person{ID: 1}, []*entity.Person{{ID: 1, Name: "Aakash", Age: 25, State: "Bihar"}}, nil},
		{entity.Person{Name: "Aakash"}, []*entity.Person{{ID: 1, Name: "Aakash", Age: 25, State: "Bihar"}}, nil},
		{entity.Person{Name: "Aakash", ID: 1, State: "Bihar", Age: 25},
			[]*entity.Person{{ID: 1, Name: "Aakash", Age: 25, State: "Bihar"}}, nil},
		{entity.Person{}, []*entity.Person{
			{ID: 1, Name: "Aakash", Age: 25, State: "Bihar"}, {ID: 3, Name: "Kali", Age: 40, State: "karnataka"}}, nil},
		{entity.Person{ID: 9, State: "Bihar"}, nil, nil},
	}

	k := initializeTest(t)
	ctx := gofr.NewContext(nil, nil, k)

	var p Person

	for i, tc := range tests {
		output := p.Get(ctx, tc.input)

		if !isSubset(output, tc.output) {
			t.Errorf("Testcase[%v] Failed\tExpected %v\nGot %v\n", i, tc.output, output)
		}
	}
}

func TestCreate(t *testing.T) {
	tests := []struct {
		input  entity.Person
		output []*entity.Person
		err    error
	}{
		{entity.Person{ID: 2, Name: "himari", Age: 30, State: "bihar"}, []*entity.Person{{ID: 2, Name: "himari", Age: 30, State: "bihar"}}, nil},
	}

	k := initializeTest(t)
	ctx := gofr.NewContext(nil, nil, k)

	var p Person

	for i, tc := range tests {
		output, err := p.Create(ctx, tc.input)

		if tc.err != err {
			t.Errorf("Testcase[%v] Failed\tExpected %v\nGot %v\n", i, tc.err, err)
		}

		if !reflect.DeepEqual(tc.output, output) {
			t.Errorf("Testcase[%v] Failed\tExpected %v\nGot %v\n", i, tc.output, output)
		}
	}
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		input  entity.Person
		output []*entity.Person
		err    error
	}{
		{entity.Person{ID: 3}, []*entity.Person{{ID: 3, Name: "Kali", Age: 40, State: "karnataka"}}, nil},
		{entity.Person{ID: 3, Name: "Mahi", Age: 40, State: "Jharkhand"},
			[]*entity.Person{{ID: 3, Name: "Mahi", Age: 40, State: "Jharkhand"}}, nil},
		{entity.Person{ID: 3, Age: 30, State: "Bihar"}, []*entity.Person{{ID: 3, Name: "Mahi", Age: 30, State: "Bihar"}}, nil},
	}

	k := initializeTest(t)
	ctx := gofr.NewContext(nil, nil, k)

	var p Person

	for i, tc := range tests {
		output, err := p.Update(ctx, tc.input)

		if tc.err != err {
			t.Errorf("Testcase[%v] Failed\tExpected %v\nGot %v\n", i, tc.err, err)
		}

		if !reflect.DeepEqual(tc.output, output) {
			t.Errorf("Testcase[%v] Failed\tExpected %v\nGot %v\n", i, tc.output, output)
		}
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		input  string
		output string
		err    error
	}{
		{"3", "3", nil},
	}

	k := initializeTest(t)
	ctx := gofr.NewContext(nil, nil, k)

	var p Person

	for i, tc := range tests {
		err := p.Delete(ctx, tc.input)

		if tc.err != err {
			t.Errorf("Testcase[%v] Failed\tExpected %v\nGot %v\n", i, tc.err, err)
		}
	}
}
