package main

import (
	"bytes"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/zopsmart/gofr/pkg/datastore"
	"github.com/zopsmart/gofr/pkg/gofr/config"
	"github.com/zopsmart/gofr/pkg/gofr/request"
	"github.com/zopsmart/gofr/pkg/log"
)

func TestMain(m *testing.M) {
	logger := log.NewLogger()
	c := config.NewGoDotEnvProvider(logger, "./configs")
	port, _ := strconv.Atoi(c.Get("CASS_DB_PORT"))
	ycqlCfg := datastore.CassandraCfg{
		Hosts:    c.Get("CASS_DB_HOST"),
		Port:     port,
		Username: c.Get("CASS_DB_USER"),
		Password: c.Get("CASS_DB_PASS"),
		Keyspace: "system",
	}

	ycqlDB, err := datastore.GetNewYCQL(logger, &ycqlCfg)
	if err != nil {
		logger.Errorf("Failed, unable to connect to ycql")
	}

	err = ycqlDB.Session.Query(
		"CREATE KEYSPACE IF NOT EXISTS test WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor': '1'} " +
			"AND DURABLE_WRITES = true;").Exec()
	if err != nil {
		logger.Errorf("Failed to create keyspace :%v", err)
	}

	ycqlCfg.Keyspace = "test"

	ycqlDB, err = datastore.GetNewYCQL(logger, &ycqlCfg)
	if err != nil {
		logger.Errorf("Failed to connect with ycql :%v", err)
	}

	// remove table if exist
	_ = ycqlDB.Session.Query("DROP TABLE IF EXISTS shop").Exec()

	queryStr := "CREATE TABLE shop (id int PRIMARY KEY, name varchar, location varchar , state varchar ) " +
		"WITH transactions = { 'enabled' : true };"

	err = ycqlDB.Session.Query(queryStr).Exec()
	if err != nil {
		logger.Errorf("Failed creation of Table shop :%v", err)
	} else {
		logger.Info("Table shop created Successfully")
	}

	os.Exit(m.Run())
}

func TestIntegrationShop(t *testing.T) {
	// call  the main function
	go main()

	time.Sleep(time.Second * 5)

	testcases := []struct {
		method             string
		endpoint           string
		expectedStatusCode int
		body               []byte
	}{
		{"GET", "shop?name=Vikash", 200, nil},
		{"POST", "shop", 201, []byte(`{"id":    4, "name":  "Puma", "location":  "Belandur" , "state": "karnataka"}`)},
		{"POST", "shop", 201, []byte(`{"id":    7, "name":  "Kalash", "location": "Jehanabad", "state": "Bihar"}`)},
		{"GET", "unknown", 404, nil},
		{"GET", "shop/id", 404, nil},
		{"PUT", "shop", 404, nil},
		{"DELETE", "shop/4", 204, nil},
	}
	for i, tc := range testcases {
		req, _ := request.NewMock(tc.method, "http://localhost:9005/"+tc.endpoint, bytes.NewBuffer(tc.body))

		cl := http.Client{}
		resp, _ := cl.Do(req)

		if resp != nil && resp.StatusCode != tc.expectedStatusCode {
			t.Errorf("Testcase[%v] Failed.\tExpected %v\tGot %v\n", i, tc.expectedStatusCode, resp.StatusCode)
		}

		resp.Body.Close()
	}
}
