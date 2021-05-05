package main

import (
	"github.com/zopsmart/gofr/examples/sample-validation/handler"
	"github.com/zopsmart/gofr/pkg/gofr"
)

func main() {
	k := gofr.New()
	k.POST("/phone", handler.ValidateEntry)

	k.Server.HTTP.Port = 9010
	k.Start()
}
