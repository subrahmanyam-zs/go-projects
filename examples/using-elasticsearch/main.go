package main

import (
	"developer.zopsmart.com/go/gofr/examples/using-elasticsearch/handler"
	"developer.zopsmart.com/go/gofr/examples/using-elasticsearch/store/customer"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	k := gofr.New()
	k.Server.ValidateHeaders = false

	h := handler.New(customer.Customer{})

	k.REST("customer", h)
	k.Server.HTTP.Port = 8001
	k.Start()
}
