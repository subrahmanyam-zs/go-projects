package main

import (
	"developer.zopsmart.com/go/gofr/examples/sample-validation/handler"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	k := gofr.New()
	k.POST("/phone", handler.ValidateEntry)

	k.Server.HTTP.Port = 9010
	k.Start()
}
