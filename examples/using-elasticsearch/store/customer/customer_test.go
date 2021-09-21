package customer

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/examples/using-elasticsearch/model"
	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
	"developer.zopsmart.com/go/gofr/pkg/log"
)

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

func initializeElasticsearchClient() (*Customer, *gofr.Context) {
	var c Customer

	k := gofr.New()
	req, _ := http.NewRequest(http.MethodGet, "/customers/_search", nil)
	r := request.NewHTTPRequest(req)
	context := gofr.NewContext(nil, r, k)
	context.Context = req.Context()

	return &c, context
}

func TestCustomer_Get(t *testing.T) {
	testCases := []struct {
		name   string
		output []model.Customer
	}{
		{"Henry", []model.Customer{{ID: "1", Name: "Henry", City: "Bangalore"}}},
		{"Random", nil},
		{"", []model.Customer{
			{ID: "1", Name: "Henry", City: "Bangalore"},
			{ID: "2", Name: "Bitsy", City: "Mysore"},
			{ID: "3", Name: "Magic", City: "Bangalore"},
		}},
	}

	for _, tc := range testCases {
		store, context := initializeElasticsearchClient()

		output, err := store.Get(context, tc.name)

		assert.Equal(t, nil, err)

		assert.Equal(t, tc.output, output)
	}
}

func TestCustomer_GetByID(t *testing.T) {
	testCases := []struct {
		id     string
		err    error
		output model.Customer
	}{
		{"1", nil, model.Customer{ID: "1", Name: "Henry", City: "Bangalore"}},
		{"", errors.EntityNotFound{Entity: "customer", ID: ""}, model.Customer{}},
	}

	for _, tc := range testCases {
		store, context := initializeElasticsearchClient()

		output, err := store.GetByID(context, tc.id)

		assert.Equal(t, tc.err, err)

		assert.Equal(t, tc.output, output)
	}
}

func TestCustomer_Create(t *testing.T) {
	var (
		input, expOutput model.Customer
	)

	input = model.Customer{ID: "4", Name: "Elon", City: "Chandigarh"}
	expOutput = model.Customer{ID: "4", Name: "Elon", City: "Chandigarh"}

	store, context := initializeElasticsearchClient()
	output, err := store.Create(context, input)

	assert.Equal(t, nil, err)

	assert.Equal(t, expOutput, output)
}

func TestCustomer_Update(t *testing.T) {
	testCases := struct {
		id     string
		input  model.Customer
		err    error
		output model.Customer
	}{
		"4", model.Customer{ID: "4", Name: "Elon", City: "Bangalore"}, nil, model.Customer{ID: "4", Name: "Elon", City: "Bangalore"},
	}

	store, context := initializeElasticsearchClient()
	output, err := store.Update(context, testCases.input, testCases.id)

	assert.Equal(t, testCases.err, err)

	assert.Equal(t, testCases.output, output)
}

func TestCustomer_Delete(t *testing.T) {
	store, context := initializeElasticsearchClient()

	err := store.Delete(context, "1")

	assert.Equal(t, nil, err)
}
