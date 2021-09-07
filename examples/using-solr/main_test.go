package main

import (
	"bytes"
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
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
		{http.MethodGet, "unknown", http.StatusNotFound, nil},
		{http.MethodGet, "/customer/id", http.StatusNotFound, nil},
		{http.MethodGet, "customer?id=2", http.StatusOK, nil},
	}

	for _, tc := range testcases {
		req, _ := request.NewMock(tc.method, "http://localhost:9099/"+tc.endpoint, bytes.NewBuffer(tc.body))
		c := http.Client{}

		resp, err := c.Do(req)
		if resp == nil || err != nil {
			t.Error(err)
			continue
		}

		if resp.StatusCode != tc.expectedStatusCode {
			t.Errorf("Failed.\tExpected %v\tGot %v\n", tc.expectedStatusCode, resp.StatusCode)
		}

		_ = resp.Body.Close()
	}
}

func TestMain(m *testing.M) {
	k := gofr.New()

	host := os.Getenv("SOLR_HOST")
	port := os.Getenv("SOLR_PORT")

	resp, err := http.Get("http://localhost:" + port + "/solr/admin/collections?action=CREATE&name=customer&numShards=1&replicationFactor=1")
	if err != nil {
		k.Logger.Errorf("error in sending request")
		os.Exit(1)
	}

	_ = resp.Body.Close()

	client := datastore.NewSolrClient(host, port)
	body := []byte(`{
	"add-field": {
		"name": "id",
        "type": "int",
         "stored": "false",
	}}`)

	document := bytes.NewBuffer(body)
	_, _ = client.AddField(context.TODO(), "customer", document)

	body = []byte(`{
		"add-field": {
			"name": "name",
				"type": "string",
				"stored": "true"
		}
	}`)

	document = bytes.NewBuffer(body)
	_, _ = client.AddField(context.TODO(), "customer", document)

	body = []byte(`{
		"add-field":{
		   "name":"dateOfBirth",
		   "type":"string",
		"stored":true }}`)

	document = bytes.NewBuffer(body)
	_, _ = client.UpdateField(context.TODO(), "customer", document)

	body = []byte(`{
		     "add-field":{
			   "name":"name",
			   "type":"string",
		    "stored":true }
			}`)
	document = bytes.NewBuffer(body)
	_, _ = client.AddField(context.TODO(), "customer", document)

	os.Exit(m.Run())
}
