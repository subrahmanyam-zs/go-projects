package customer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"developer.zopsmart.com/go/gofr/examples/using-elasticsearch/model"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Customer struct{}

const index = "customers"

func (c Customer) Get(context *gofr.Context, name string) ([]model.Customer, error) {
	var body string

	if name != "" {
		body = fmt.Sprintf(`{"query" : { "match" : {"name":"%s"} }}`, name)
	}

	es := context.Elasticsearch

	res, err := es.Search(
		es.Search.WithIndex(index),
		es.Search.WithContext(context),
		es.Search.WithBody(strings.NewReader(body)),
		es.Search.WithPretty(),
	)
	if err != nil {
		return nil, errors.DB{Err: err}
	}

	var customers []model.Customer

	err = es.BindArray(res, &customers)
	if err != nil {
		return nil, err
	}

	return customers, nil
}

func (c Customer) GetByID(context *gofr.Context, id string) (model.Customer, error) {
	var customer model.Customer

	es := context.Elasticsearch

	res, err := es.Search(
		es.Search.WithIndex(index),
		es.Search.WithContext(context),
		es.Search.WithBody(strings.NewReader(fmt.Sprintf(`{"query" : { "match" : {"id":"%s"} }}`, id))),
		es.Search.WithPretty(),
		es.Search.WithSize(1),
	)
	if err != nil {
		return customer, errors.DB{Err: err}
	}

	err = es.Bind(res, &customer)
	if err != nil {
		return customer, err
	}

	if customer.ID == "" {
		return customer, errors.EntityNotFound{Entity: "customer", ID: id}
	}

	return customer, nil
}

func (c Customer) Update(context *gofr.Context, customer model.Customer, id string) (model.Customer, error) {
	body, err := json.Marshal(customer)
	if err != nil {
		return model.Customer{}, errors.DB{Err: err}
	}

	es := context.Elasticsearch

	res, err := es.Index(
		index,
		bytes.NewReader(body),
		es.Index.WithRefresh("true"),
		es.Index.WithPretty(),
		es.Index.WithContext(context),
		es.Index.WithDocumentID(id),
	)
	if err != nil {
		return model.Customer{}, errors.DB{Err: err}
	}

	resp, err := es.Body(res)
	if err != nil {
		return model.Customer{}, errors.DB{Err: err}
	}

	if id, ok := resp["_id"].(string); ok {
		return c.GetByID(context, id)
	}

	return model.Customer{}, errors.Error("update error: invalid id")
}

func (c Customer) Create(context *gofr.Context, customer model.Customer) (model.Customer, error) {
	body, err := json.Marshal(customer)
	if err != nil {
		return model.Customer{}, errors.DB{Err: err}
	}

	es := context.Elasticsearch

	res, err := es.Index(
		index,
		bytes.NewReader(body),
		es.Index.WithRefresh("true"),
		es.Index.WithPretty(),
		es.Index.WithContext(context),
		es.Index.WithDocumentID(customer.ID),
	)
	if err != nil {
		return model.Customer{}, errors.DB{Err: err}
	}

	resp, err := es.Body(res)
	if err != nil {
		return model.Customer{}, errors.DB{Err: err}
	}

	if id, ok := resp["_id"].(string); ok {
		return c.GetByID(context, id)
	}

	return model.Customer{}, errors.Error("create error: invalid id")
}

func (c Customer) Delete(context *gofr.Context, id string) error {
	es := context.Elasticsearch

	_, err := es.Delete(
		index,
		id,
		es.Delete.WithContext(context),
		es.Delete.WithPretty(),
	)
	if err != nil {
		return errors.DB{Err: err}
	}

	return nil
}
