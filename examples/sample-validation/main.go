package main

import (
	"developer.zopsmart.com/go/gofr/examples/sample-validation/handler"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	app := gofr.New()
	app.POST("/phone", handler.ValidateEntry)

	app.Server.HTTP.Port = 9010
	app.Start()
}
