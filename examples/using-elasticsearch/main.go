package main

import (
	"github.com/zopsmart/gofr/examples/using-elasticsearch/handler"
	"github.com/zopsmart/gofr/examples/using-elasticsearch/store/customer"
	"github.com/zopsmart/gofr/pkg/gofr"
)

func main() {
	k := gofr.New()
	k.Server.ValidateHeaders = false

	h := handler.New(customer.Customer{})

	k.REST("customer", h)
	k.Server.HTTP.Port = 8001
	k.Start()
}
