package main

import (
	"bytes"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/zopsmart/gofr/pkg/gofr"
	"github.com/zopsmart/gofr/pkg/gofr/request"
)

func TestMain(m *testing.M) {
	k := gofr.New()

	db := k.DB()
	if db == nil {
		return
	}

	query := `
 	   CREATE TABLE IF NOT EXISTS customers (
	   id serial primary key,
	   name varchar (50))
	`

	if k.Config.Get("DB_DIALECT") == "mssql" {
		query = `
		IF NOT EXISTS
	(  SELECT [name]
		FROM sys.tables
      WHERE [name] = 'customers'
   ) CREATE TABLE customers (id int primary key identity(1,1),
	   name varchar (50))
	`
	}

	if _, err := db.Exec(query); err != nil {
		k.Logger.Errorf("got error sourcing the schema: ", err)
	}

	os.Exit(m.Run())
}

func TestIntegration(t *testing.T) {
	go main()
	time.Sleep(time.Second * 5)

	tcs := []struct {
		method             string
		endpoint           string
		expectedStatusCode int
		body               []byte
	}{
		{"GET", "customer", 200, nil},
	}

	for _, tc := range tcs {
		req, _ := request.NewMock(tc.method, "http://localhost:9092/"+tc.endpoint, bytes.NewBuffer(tc.body))
		c := http.Client{}

		//nolint: bodyclose, no response body to close
		resp, _ := c.Do(req)
		if resp != nil && resp.StatusCode != tc.expectedStatusCode {
			t.Errorf("Failed.\tExpected %v\tGot %v\n", tc.expectedStatusCode, resp.StatusCode)
		}
	}
}
