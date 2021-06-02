package main

import (
	"bytes"
	"net/http"
	"os"
	"testing"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
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

func TestIntegrationPersons(t *testing.T) {
	// call  the main function
	go main()

	time.Sleep(time.Second * 5)

	testcases := []struct {
		method             string
		endpoint           string
		expectedStatusCode int
		body               []byte
	}{
		{"GET", "persons?name=Vikash", 200, nil},
		{"POST", "persons", 201, []byte(`{"id":    7, "name":  "Kali", "age":   40, "State": "karnataka"}`)},
		{"POST", "persons", 201, []byte(`{"id":    8, "name":  "Kali"}`)},
		{"DELETE", "persons/7", 204, nil},
		{"GET", "unknown", 404, nil},
		{"GET", "persons/id", 404, nil},
		{"PUT", "persons", 404, nil},
	}
	for i, tc := range testcases {
		req, _ := request.NewMock(tc.method, "http://localhost:9094/"+tc.endpoint, bytes.NewBuffer(tc.body))

		cl := http.Client{}
		resp, _ := cl.Do(req)

		if resp != nil && resp.StatusCode != tc.expectedStatusCode {
			t.Errorf("Testcase[%v] Failed.\tExpected %v\tGot %v\n", i, tc.expectedStatusCode, resp.StatusCode)
		}

		resp.Body.Close()
	}
}
