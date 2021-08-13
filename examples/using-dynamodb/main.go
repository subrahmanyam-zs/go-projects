package main

import (
	"developer.zopsmart.com/go/gofr/examples/using-dynamodb/handlers"
	"developer.zopsmart.com/go/gofr/examples/using-dynamodb/store/person"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	k := gofr.New()

	store := person.New("person")
	handler := handlers.New(store)

	k.GET("/person/{id}", handler.GetByID)
	k.POST("/person", handler.Create)
	k.PUT("/person/{id}", handler.Update)
	k.DELETE("/person/{id}", handler.Delete)

	k.Start()
}
