package main

import (
	"developer.zopsmart.com/go/gofr/examples/using-http-service/handler"
	svc "developer.zopsmart.com/go/gofr/examples/using-http-service/service"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	k := gofr.New()

	catSvc := svc.New("http://catalog-service/brand", k.Logger)
	h := handler.New(catSvc)

	k.GET("/brand/{id}", h.Get)
	k.Start()
}
