package shop

import (
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/zopsmart/gofr/examples/using-ycql/entity"
	"github.com/zopsmart/gofr/pkg/datastore"
	"github.com/zopsmart/gofr/pkg/gofr"
	"github.com/zopsmart/gofr/pkg/gofr/config"
	"github.com/zopsmart/gofr/pkg/log"
)

func TestMain(m *testing.M) {
	logger := log.NewLogger()
	c := config.NewGoDotEnvProvider(logger, "../../configs")
	cassandraPort, _ := strconv.Atoi(c.Get("CASS_DB_PORT"))
	ycqlCfg := datastore.CassandraCfg{
		Hosts:    c.Get("CASS_DB_HOST"),
		Port:     cassandraPort,
		Username: c.Get("CASS_DB_USER"),
		Password: c.Get("CASS_DB_PASS"),
		Keyspace: "system",
	}

	ycqlDB, err := datastore.GetNewYCQL(logger, &ycqlCfg)
	if err != nil {
		logger.Errorf("got error while connection")
	}

	err = ycqlDB.Session.Query(
		"CREATE KEYSPACE IF NOT EXISTS test WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor': '1'} " +
			"AND DURABLE_WRITES = true;").Exec()
	if err != nil {
		logger.Errorf("got error while connection")
	}

	ycqlCfg.Keyspace = "test"

	ycqlDB, err = datastore.GetNewYCQL(logger, &ycqlCfg)
	if err != nil {
		logger.Errorf("got error while connection")
	}

	os.Exit(m.Run())
}

func initializeTest(t *testing.T) *gofr.Context {
	k := gofr.New()
	// initializing the seeder
	seeder := datastore.NewSeeder(&k.DataStore, "../../db")
	query := "CREATE TABLE IF NOT EXISTS shop (id int PRIMARY KEY, name varchar, location varchar , state varchar ) " +
		"WITH transactions = { 'enabled' : true };"

	err := k.YCQL.Session.Query(query).Exec()
	if err != nil {
		t.Errorf("Table shop is not created")
	}

	ctx := gofr.NewContext(nil, nil, k)

	seeder.RefreshYCQL(t, "shop")

	return ctx
}

func TestGet(t *testing.T) {
	tests := []struct {
		input  entity.Shop
		output []entity.Shop
		err    error
	}{
		{entity.Shop{ID: 1}, []entity.Shop{{ID: 1, Name: "Pramod", Location: "Gaya", State: "Bihar"}}, nil},
		{entity.Shop{Name: "Pramod"}, []entity.Shop{{ID: 1, Name: "Pramod", Location: "Gaya", State: "Bihar"}}, nil},
		{entity.Shop{Name: "Pramod", ID: 1, State: "Bihar", Location: "Gaya"},
			[]entity.Shop{{ID: 1, Name: "Pramod", Location: "Gaya", State: "Bihar"}}, nil},
		{entity.Shop{}, []entity.Shop{
			{ID: 1, Name: "Pramod", Location: "Gaya", State: "Bihar"}, {ID: 2, Name: "Shubh", Location: "HSR", State: "Karnataka"}}, nil},
		{entity.Shop{ID: 9, State: "Bihar"}, nil, nil},
	}

	ctx := initializeTest(t)

	var p Shop

	for i, tc := range tests {
		output := p.Get(ctx, tc.input)

		if !reflect.DeepEqual(tc.output, output) {
			t.Errorf("Testcase[%v] Failed\tExpected %v\nGot %v\n", i, tc.output, output)
		}
	}
}

func TestCreate(t *testing.T) {
	tests := []struct {
		input  entity.Shop
		output []entity.Shop
		err    error
	}{
		{entity.Shop{ID: 2, Name: "himalaya", Location: "Gaya", State: "bihar"},
			[]entity.Shop{{ID: 2, Name: "himalaya", Location: "Gaya", State: "bihar"}}, nil},
	}

	ctx := initializeTest(t)

	var p Shop

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
		input  entity.Shop
		output []entity.Shop
		err    error
	}{
		{entity.Shop{ID: 2}, []entity.Shop{{ID: 2, Name: "Shubh", Location: "HSR", State: "Karnataka"}}, nil},
		{entity.Shop{ID: 2, Name: "Mahi", Location: "Dhanbad", State: "Jharkhand"},
			[]entity.Shop{{ID: 2, Name: "Mahi", Location: "Dhanbad", State: "Jharkhand"}}, nil},
		{entity.Shop{ID: 2, Location: "Gaya", State: "Bihar"}, []entity.Shop{{ID: 2, Name: "Mahi", Location: "Gaya", State: "Bihar"}}, nil},
	}

	ctx := initializeTest(t)

	var p Shop

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

	ctx := initializeTest(t)

	var p Shop

	for i, tc := range tests {
		err := p.Delete(ctx, tc.input)

		if tc.err != err {
			t.Errorf("Testcase[%v] Failed\tExpected %v\nGot %v\n", i, tc.err, err)
		}
	}
}
