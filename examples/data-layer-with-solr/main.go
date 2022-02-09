package main

import (
	"os"

	"developer.zopsmart.com/go/gofr/examples/data-layer-with-solr/handler"
	"developer.zopsmart.com/go/gofr/examples/data-layer-with-solr/store/customer"
	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	app := gofr.New()

	// initializing the solr client for core layer
	client := datastore.NewSolrClient(os.Getenv("SOLR_HOST"), os.Getenv("SOLR_PORT"))
	customerCore := customer.New(client)
	customerConsumer := handler.New(customerCore)

	// Specifying the different routes supported by this service
	app.GET("/customer", customerConsumer.List)
	app.POST("/customer", customerConsumer.Create)
	app.PUT("/customer", customerConsumer.Update)
	app.DELETE("/customer", customerConsumer.Delete)

	app.Start()
}
