package customer

import (
	"errors"
	"net/http"
	"testing"

	"github.com/olivere/elastic/v6"
	"developer.zopsmart.com/go/gofr/examples/using-elasticsearch/model"
	errors2 "developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
)

func initializeElasticsearchClient() (*Customer, *gofr.Context) {
	var c Customer

	k := gofr.New()
	req, _ := http.NewRequest("GET", "/customer/1", nil)
	r := request.NewHTTPRequest(req)
	context := gofr.NewContext(nil, r, k)
	context.Context = req.Context()

	return &c, context
}

func TestCustomer_Get(t *testing.T) {
	testCases := []struct {
		name string
		err  error
	}{
		{"", nil},
		{"Miracle", nil},
	}

	for i, v := range testCases {
		m, context := initializeElasticsearchClient()
		_, _ = m.Create(context, model.Customer{Name: "Miracle"})

		_, err := m.Get(context, v.name)

		if err != v.err {
			t.Errorf("[TESTCASE%d]Failed.Got: %v\tExpected: %v\n", i+1, err, v.err)
		}
	}
}

func TestCustomer_GetByIDMissingId(t *testing.T) {
	expectedErr := errors.New("missing required fields: [Id]")
	m, context := initializeElasticsearchClient()
	_, err := m.GetByID(context, "")

	if err.Error() != expectedErr.Error() {
		t.Errorf("Failed.Got: %v\tExpected: %v\n", err, expectedErr)
	}
}

func TestCustomer_GetByIDNotFound(t *testing.T) {
	expectedError := elastic.Error{Status: 404, Details: nil}
	m, context := initializeElasticsearchClient()
	_, err := m.GetByID(context, "RpxS9W8BvEf544")
	dbErr, _ := err.(errors2.DB)
	e, ok := dbErr.Err.(*elastic.Error)

	if !ok && e.Status != expectedError.Status {
		t.Errorf("Failed.Got: %v\tExpected: %v\n", err, expectedError)
	}
}

func TestCustomer_GetByID(t *testing.T) {
	m, context := initializeElasticsearchClient()
	_, err := m.GetByID(context, customerID)

	if err != nil {
		t.Errorf("Expected nil in GetById, Got %v", err)
	}
}

func TestCustomer_Create(t *testing.T) {
	testCases := []struct {
		customer model.Customer
		err      error
	}{
		{model.Customer{}, nil},
		{model.Customer{Name: "Eron"}, nil},
	}

	for i, v := range testCases {
		m, context := initializeElasticsearchClient()
		_, err := m.Create(context, v.customer)

		if err != v.err {
			t.Errorf("[TESTCASE%d]Failed.Got: %v\tExpected: %v\n", i+1, err, v.err)
		}
	}
}

func TestCustomer_Update(t *testing.T) {
	testCases := []struct {
		customer model.Customer
		id       string
		err      error
	}{
		{model.Customer{Name: "Heu"}, customerID, nil},
		{model.Customer{Name: "Marc"}, customerID, nil},
	}

	for i, v := range testCases {
		m, context := initializeElasticsearchClient()
		_, err := m.Update(context, v.customer, v.id)

		if err != v.err {
			t.Errorf("[TESTCASE%d]Failed.Got: %v\tExpected: %v\n", i+1, err, v.err)
		}
	}
}

func TestCustomer_UpdateError(t *testing.T) {
	testCases := []struct {
		customer model.Customer
		id       string
		err      elastic.Error
	}{
		{model.Customer{Name: "Magic"}, "RpxS9W8BvEf-ncpuarrhwK", elastic.Error{Status: 404}},
		{model.Customer{Name: "Eron"}, "", elastic.Error{Status: 400}},
	}

	for i, v := range testCases {
		m, context := initializeElasticsearchClient()
		_, err := m.Update(context, v.customer, v.id)
		dbErr, _ := err.(errors2.DB)
		e, ok := dbErr.Err.(*elastic.Error)

		if !ok && e.Status != v.err.Status {
			t.Errorf("[TESTCASE %d]Failed.Got: %v\tExpected: %v\n", i+1, err, v.err)
		}
	}
}

func TestCustomer_DeleteError(t *testing.T) {
	expectedErr := elastic.Error{Status: 404}
	m, context := initializeElasticsearchClient()

	err := m.Delete(context, "yu6545343")
	dbErr, _ := err.(errors2.DB)
	e, ok := dbErr.Err.(*elastic.Error)

	if !ok || e.Status != expectedErr.Status {
		t.Errorf("Failed.Got: %v\tExpected: %v\n", err, expectedErr)
	}
}

func TestCustomer_Delete(t *testing.T) {
	m, context := initializeElasticsearchClient()
	err := m.Delete(context, customerID)

	if err != nil {
		t.Errorf("Expected successful delete but got: %v", err)
	}
}
