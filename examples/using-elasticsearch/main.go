package main

import (
	"developer.zopsmart.com/go/gofr/examples/using-elasticsearch/handler"
	"developer.zopsmart.com/go/gofr/examples/using-elasticsearch/store/customer"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	k := gofr.New()

	h := handler.New(customer.Customer{})

	k.REST("customer", h)

	k.Start()
}
