package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"

	"developer.zopsmart.com/go/gofr/examples/using-elasticsearch/model"
	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/log"
)

const index = "customers"

// creating the index 'customers' and populating data to use it in tests
func TestMain(m *testing.M) {
	k := gofr.New()

	const mapping = `{"settings": {
		"number_of_shards": 1},
	"mappings": {
		"_doc": {
			"properties": {
				"id": {"type": "text"},
				"name": {"type": "text"},
				"city": {"type": "text"}
			}}}}`

	client := k.Elasticsearch

	_, err := client.Indices.Delete([]string{index}, client.Indices.Delete.WithIgnoreUnavailable(true))
	if err != nil {
		k.Logger.Errorf("error deleting index: %s", err.Error())
	}

	_, err = client.Indices.Create(index,
		client.Indices.Create.WithBody(bytes.NewReader([]byte(mapping))),
		client.Indices.Create.WithPretty(),
	)
	if err != nil {
		k.Logger.Errorf("error creating index: %s", err.Error())
	}

	insert(k.Logger, model.Customer{ID: "1", Name: "Henry", City: "Bangalore"}, client)
	insert(k.Logger, model.Customer{ID: "2", Name: "Bitsy", City: "Mysore"}, client)
	insert(k.Logger, model.Customer{ID: "3", Name: "Magic", City: "Bangalore"}, client)

	os.Exit(m.Run())
}

func insert(logger log.Logger, customer model.Customer, client datastore.Elasticsearch) {
	body, _ := json.Marshal(customer)

	_, err := client.Index(
		index,
		bytes.NewReader(body),
		client.Index.WithRefresh("true"),
		client.Index.WithPretty(),
		client.Index.WithDocumentID(customer.ID),
	)
	if err != nil {
		logger.Errorf("error inserting documents: %s", err.Error())
	}
}

func TestRoutes(t *testing.T) {
	go main()
	time.Sleep(time.Second * 2)

	testcases := []struct {
		method             string
		endpoint           string
		expectedStatusCode int
		body               []byte
	}{
		{http.MethodGet, "customer", http.StatusOK, nil},
		{http.MethodGet, "customer/7", http.StatusInternalServerError, nil},
		{http.MethodPost, "customer", http.StatusCreated, []byte(`{"id":"1","name":"test","city":"xyz"}`)},
		{http.MethodPut, "customer/1", http.StatusOK, []byte(`{"id":"1","name":"test1","city":"xyz2"}`)},
		{http.MethodDelete, "customer/1", http.StatusNoContent, nil},
	}

	for _, tc := range testcases {
		req, _ := http.NewRequest(tc.method, "http://localhost:8001/"+tc.endpoint, bytes.NewBuffer(tc.body))
		c := http.Client{}

		resp, err := c.Do(req)
		if err != nil {
			t.Errorf("error while making request: %v", err)
		}

		if resp != nil && resp.StatusCode != tc.expectedStatusCode {
			t.Errorf("Failed.\tExpected %v\tGot %v\n", tc.expectedStatusCode, resp.StatusCode)
		}

		_ = resp.Body.Close()
	}
}
