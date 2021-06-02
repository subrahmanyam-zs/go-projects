package customer

import (
	"context"
	"os"
	"testing"

	"developer.zopsmart.com/go/gofr/examples/using-elasticsearch/model"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

// nolint:gochecknoglobals // needed for testing customer package
var customerID string

// creating the index 'customers' and populating data to use it in tests
func TestMain(m *testing.M) {
	k := gofr.New()

	const mapping = `{"settings": {
		"number_of_shards": 1},
	"mappings": {
		"_doc": {
			"properties": {
				"name": {"type": "text"},
				"phone": {"type": "text"},
				"city": {"type": "text"}
			}}}}`

	client := k.Elasticsearch
	exists, _ := client.IndexExists("customers").Do(context.Background())

	if !exists {
		// Create a new index.
		_, _ = client.CreateIndex("customers").BodyString(mapping).Do(context.Background())
	} else {
		_, _ = client.DeleteIndex("customers").Do(context.Background())
		_, _ = client.CreateIndex("customers").BodyString(mapping).Do(context.Background())
	}

	customer := model.Customer{Name: "Henry", City: "Bangalore"}
	res, _ := client.Index().Index("customers").Type("_doc").BodyJson(customer).Do(context.Background())

	customerID = res.Id

	customer = model.Customer{Name: "Bitsy", City: "Mysore"}
	_, _ = client.Index().Index("customers").Type("_doc").BodyJson(customer).Do(context.Background())

	customer = model.Customer{Name: "Magic", City: "Bangalore"}
	_, _ = client.Index().Index("customers").Type("_doc").BodyJson(customer).Do(context.Background())

	os.Exit(m.Run())
}
