package handler

import (
	"Customer/models"
	"Customer/service"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type Handler struct {
	service service.Customer
}

func New(customer service.Customer) Handler {
	return Handler{service: customer}
}

func (h Handler) Post(w http.ResponseWriter, r *http.Request) {
	customer, err := getCustomer(r)
	if err != nil {

	}

	customer, err = h.service.Post(r.Context(), &customer)
}

func getCustomer(r *http.Request) (models.Customer, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return models.Customer{}, errors.New("body")
	}

	var customer models.Customer

	err = json.Unmarshal(body, &customer)
	if err != nil {
		return models.Customer{}, errors.New("error during unmarshal")
	}

	return customer, nil
}
