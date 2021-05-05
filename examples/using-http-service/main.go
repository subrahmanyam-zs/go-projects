package main

import (
	"github.com/zopsmart/gofr/examples/using-http-service/handler"
	svc "github.com/zopsmart/gofr/examples/using-http-service/service"
	"github.com/zopsmart/gofr/pkg/gofr"
)

func main() {
	k := gofr.New()

	catSvc := svc.New("http://catalog-service/brand", k.Logger)
	h := handler.New(catSvc)

	k.GET("/brand/{id}", h.Get)
	k.Start()
}
