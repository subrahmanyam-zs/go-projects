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
	app := gofr.New()

	const mapping = `{"settings": {
		"number_of_shards": 1},
	"mappings": {
		"_doc": {
			"properties": {
				"id": {"type": "text"},
				"name": {"type": "text"},
				"city": {"type": "text"}
			}}}}`

	es := app.Elasticsearch

	_, err := es.Indices.Delete([]string{index}, es.Indices.Delete.WithIgnoreUnavailable(true))
	if err != nil {
		app.Logger.Errorf("error deleting index: %s", err.Error())
	}

	_, err = es.Indices.Create(index,
		es.Indices.Create.WithBody(bytes.NewReader([]byte(mapping))),
		es.Indices.Create.WithPretty(),
	)
	if err != nil {
		app.Logger.Errorf("error creating index: %s", err.Error())
	}

	insert(app.Logger, model.Customer{ID: "1", Name: "Henry", City: "Bangalore"}, es)
	insert(app.Logger, model.Customer{ID: "2", Name: "Bitsy", City: "Mysore"}, es)
	insert(app.Logger, model.Customer{ID: "3", Name: "Magic", City: "Bangalore"}, es)

	os.Exit(m.Run())
}

func insert(logger log.Logger, customer model.Customer, es datastore.Elasticsearch) {
	body, _ := json.Marshal(customer)

	_, err := es.Index(
		index,
		bytes.NewReader(body),
		es.Index.WithRefresh("true"),
		es.Index.WithPretty(),
		es.Index.WithDocumentID(customer.ID),
	)
	if err != nil {
		logger.Errorf("error inserting documents: %s", err.Error())
	}
}

func TestRoutes(t *testing.T) {
	go main()
	time.Sleep(time.Second * 2)

	tests := []struct {
		desc       string
		method     string
		endpoint   string
		statusCode int
		body       []byte
	}{
		{"get all customer success case", http.MethodGet, "customer", http.StatusOK, nil},
		{"get non existent customer", http.MethodGet, "customer/7", http.StatusInternalServerError, nil},
		{"create success", http.MethodPost, "customer", http.StatusCreated, []byte(`{"id":"1","name":"test","city":"xyz"}`)},
		{"update success", http.MethodPut, "customer/1", http.StatusOK, []byte(`{"id":"1","name":"test1","city":"xyz2"}`)},
		{"delete success", http.MethodDelete, "customer/1", http.StatusNoContent, nil},
	}

	for i, tc := range tests {
		req, _ := http.NewRequest(tc.method, "http://localhost:8001/"+tc.endpoint, bytes.NewBuffer(tc.body))
		c := http.Client{}

		resp, err := c.Do(req)
		if err != nil {
			t.Errorf("TEST[%v] Failed.\tHTTP request encountered Err: %v\n%s", i, err, tc.desc)
			continue
		}

		if resp.StatusCode != tc.statusCode {
			t.Errorf("TEST[%v] Failed.\tExpected %v\tGot %v\n%s", i, tc.statusCode, resp.StatusCode, tc.desc)
		}

		_ = resp.Body.Close()
	}
}
