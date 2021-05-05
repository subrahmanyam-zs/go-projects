package main

import (
	"bytes"
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/zopsmart/gofr/pkg/datastore"
	"github.com/zopsmart/gofr/pkg/gofr"
	"github.com/zopsmart/gofr/pkg/gofr/request"
)

func TestRoutes(t *testing.T) {
	go main()

	time.Sleep(time.Second * 5)

	testcases := []struct {
		method             string
		endpoint           string
		expectedStatusCode int
		body               []byte
	}{
		{"GET", "unknown", 404, nil},
		{"GET", "/customer/id", 404, nil},
		{"GET", "customer?id=2", 500, nil},
	}

	for _, tc := range testcases {
		req, _ := request.NewMock(tc.method, "http://localhost:9099/"+tc.endpoint, bytes.NewBuffer(tc.body))
		c := http.Client{}

		resp, _ := c.Do(req)

		if resp != nil && resp.StatusCode != tc.expectedStatusCode {
			t.Errorf("Failed.\tExpected %v\tGot %v\n", tc.expectedStatusCode, resp.StatusCode)
		}
		resp.Body.Close()
	}
}

func TestMain(m *testing.M) {
	gofr.New()

	host := os.Getenv("SOLR_HOST")
	port := os.Getenv("SOLR_PORT")
	//nolint:bodyclose //response body must be closed
	_, _ = http.Get("http://localhost:2020/solr/admin/collections?action=CREATE&name=customer&numShards=1&replicationFactor=1")

	client := datastore.NewSolrClient(host, port)
	body := []byte(`{
	"add-field": {
		"name": "id",
        "type": "int",
         "stored": "false",
	}}`)

	document := bytes.NewBuffer(body)
	_, _ = client.AddField(context.TODO(), "customers", document)

	body = []byte(`{
		"add-field": {
			"name": "name",
				"type": "string",
				"stored": "true"
		}
	}`)

	document = bytes.NewBuffer(body)
	_, _ = client.AddField(context.TODO(), "customers", document)

	body = []byte(`{
		"add-field":{
		   "name":"dateOfBirth",
		   "type":"string",
		"stored":true }}`)

	document = bytes.NewBuffer(body)
	_, _ = client.UpdateField(context.TODO(), "customers", document)

	body = []byte(`{
		     "add-field":{
			   "name":"name",
			   "type":"string",
		    "stored":true }
			}`)
	document = bytes.NewBuffer(body)
	_, _ = client.AddField(context.TODO(), "customers", document)

	os.Exit(m.Run())
}
