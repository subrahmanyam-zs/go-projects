package main

import (
	"developer.zopsmart.com/go/gofr/examples/using-elasticsearch/handler"
	"developer.zopsmart.com/go/gofr/examples/using-elasticsearch/store/customer"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	app := gofr.New()

	app.Server.ValidateHeaders = false

	store := customer.New()
	h := handler.New(store)

	app.REST("customer", h)

	app.Start()
}
