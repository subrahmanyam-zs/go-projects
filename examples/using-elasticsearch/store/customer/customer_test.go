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

func initializeElasticsearchClient() (*store, *gofr.Context) {
	store := New()

	app := gofr.New()
	req, _ := http.NewRequest(http.MethodGet, "/customers/_search", nil)
	r := request.NewHTTPRequest(req)
	ctx := gofr.NewContext(nil, r, app)
	ctx.Context = req.Context()

	return &store, ctx
}

func TestCustomer_Get(t *testing.T) {
	tests := []struct {
		desc string
		name string
		resp []model.Customer
	}{
		{"get success", "Henry", []model.Customer{{ID: "1", Name: "Henry", City: "Bangalore"}}},
		{"get non existent entity", "Random", nil},
		{"get multiple entities", "", []model.Customer{
			{ID: "1", Name: "Henry", City: "Bangalore"},
			{ID: "2", Name: "Bitsy", City: "Mysore"},
			{ID: "3", Name: "Magic", City: "Bangalore"},
		}},
	}

	for i, tc := range tests {
		store, ctx := initializeElasticsearchClient()

		resp, err := store.Get(ctx, tc.name)

		assert.Equal(t, nil, err, "TEST[%d], failed.\n%s", i, tc.desc)

		assert.Equal(t, tc.resp, resp, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

func TestCustomer_GetByID(t *testing.T) {
	tests := []struct {
		desc string
		id   string
		err  error
		resp model.Customer
	}{
		{"get by id success", "1", nil, model.Customer{ID: "1", Name: "Henry", City: "Bangalore"}},
		{"get by id fail", "", errors.EntityNotFound{Entity: "customer", ID: ""}, model.Customer{}},
	}

	for i, tc := range tests {
		store, ctx := initializeElasticsearchClient()

		resp, err := store.GetByID(ctx, tc.id)

		assert.Equal(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)

		assert.Equal(t, tc.resp, resp, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

func TestCustomer_Create(t *testing.T) {
	var (
		input, expResp model.Customer
	)

	input = model.Customer{ID: "4", Name: "Elon", City: "Chandigarh"}
	expResp = model.Customer{ID: "4", Name: "Elon", City: "Chandigarh"}

	store, ctx := initializeElasticsearchClient()
	resp, err := store.Create(ctx, input)

	assert.Equal(t, nil, err)

	assert.Equal(t, expResp, resp)
}

func TestCustomer_Update(t *testing.T) {
	tests := []struct {
		desc  string
		id    string
		input model.Customer
		err   error
		resp  model.Customer
	}{
		{"update existent entity", "4", model.Customer{ID: "4", Name: "Elon", City: "Bangalore"}, nil,
			model.Customer{ID: "4", Name: "Elon", City: "Bangalore"}},
		{"update non existent entity", "444", model.Customer{ID: "444", Name: "Musk", City: "Bangalore"}, nil,
			model.Customer{ID: "444", Name: "Musk", City: "Bangalore"}},
	}

	for i, tc := range tests {
		store, ctx := initializeElasticsearchClient()
		resp, err := store.Update(ctx, tc.input, tc.id)

		assert.Equal(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)

		assert.Equal(t, tc.resp, resp, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

func TestCustomer_Delete(t *testing.T) {
	store, ctx := initializeElasticsearchClient()

	err := store.Delete(ctx, "1")

	assert.Equal(t, nil, err)
}
