package main

import (
	handlers "developer.zopsmart.com/go/gofr/examples/data-layer-with-dynamodb/handlers/person"
	stores "developer.zopsmart.com/go/gofr/examples/data-layer-with-dynamodb/stores/person"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	app := gofr.New()

	s := stores.New("person")
	h := handlers.New(s)

	app.GET("/person/{id}", h.GetByID)
	app.POST("/person", h.Create)
	app.PUT("/person/{id}", h.Update)
	app.DELETE("/person/{id}", h.Delete)

	app.Start()
}
