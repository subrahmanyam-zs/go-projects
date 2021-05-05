package main

import (
	"os"

	"github.com/zopsmart/gofr/examples/using-solr/handler"
	"github.com/zopsmart/gofr/examples/using-solr/store/customer"
	"github.com/zopsmart/gofr/pkg/datastore"
	"github.com/zopsmart/gofr/pkg/gofr"
)

func main() {
	k := gofr.New()

	// initializing the solr client for core layer
	client := datastore.NewSolrClient(os.Getenv("SOLR_HOST"), os.Getenv("SOLR_PORT"))
	customerCore := customer.New(client)
	customerConsumer := handler.New(customerCore)

	// Specifying the different routes supported by this service
	k.GET("/customer", customerConsumer.List)
	k.POST("/customer", customerConsumer.Create)
	k.PUT("/customer", customerConsumer.Update)
	k.DELETE("/customer", customerConsumer.Delete)

	k.Start()
}
