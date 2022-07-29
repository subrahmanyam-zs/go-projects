package service

import "Customer/datastore"

type Service struct {
	customer datastore.Customer
}

func New(c datastore.Customer) Service {
	return Service{c}

}
