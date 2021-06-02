package customer

import (
	"encoding/json"

	"github.com/olivere/elastic/v6"
	"developer.zopsmart.com/go/gofr/examples/using-elasticsearch/model"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Customer struct{}

func (c Customer) Get(context *gofr.Context, name string) ([]model.Customer, error) {
	var query elastic.Query
	query = elastic.MatchAllQuery{}

	if name != "" {
		query = elastic.NewMatchQuery("name", name)
	}

	searchResult, err := context.Elasticsearch.Search().Index("customers").Query(query).Pretty(true).Do(context)
	if err != nil {
		return nil, errors.DB{Err: err}
	}

	resp := make([]model.Customer, 0)

	// Iterate through results and populate customers based on search results
	for _, hit := range searchResult.Hits.Hits {
		var c model.Customer

		b, _ := hit.Source.MarshalJSON()

		_ = json.Unmarshal(b, &c)

		c.ID = hit.Id
		resp = append(resp, c)
	}

	return resp, nil
}

func (c Customer) GetByID(context *gofr.Context, id string) (*model.Customer, error) {
	var customer model.Customer

	searchResult, err := context.Elasticsearch.Get().Index("customers").Id(id).Pretty(true).Do(context)
	if err != nil {
		return nil, errors.DB{Err: err}
	}

	b, _ := searchResult.Source.MarshalJSON()

	err = json.Unmarshal(b, &customer)
	if err != nil {
		return nil, err
	}

	customer.ID = searchResult.Id

	return &customer, nil
}

func (c Customer) Update(context *gofr.Context, customer model.Customer, id string) (*model.Customer, error) {
	resp, err := context.Elasticsearch.Update().Index("customers").Type("_doc").Id(id).Doc(map[string]interface{}{
		"name": customer.Name}).Do(context)
	if err != nil {
		return nil, errors.DB{Err: err}
	}

	customer.ID = resp.Id

	return &customer, nil
}

func (c Customer) Create(context *gofr.Context, customer model.Customer) (*model.Customer, error) {
	resp, err := context.Elasticsearch.Index().Index("customers").Type("_doc").BodyJson(customer).Do(context)
	if err != nil {
		return nil, errors.DB{Err: err}
	}

	customer.ID = resp.Id

	return &customer, nil
}

func (c Customer) Delete(context *gofr.Context, id string) error {
	_, err := context.Elasticsearch.Delete().Type("_doc").Index("customers").Id(id).Do(context)
	if err != nil {
		return errors.DB{Err: err}
	}

	return err
}
