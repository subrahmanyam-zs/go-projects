package main

import (
	handlers "developer.zopsmart.com/go/gofr/examples/using-http-service/handlers/user"
	services "developer.zopsmart.com/go/gofr/examples/using-http-service/services/user"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	svc "developer.zopsmart.com/go/gofr/pkg/service"
)

func main() {
	app := gofr.New()

	sampleSvc := svc.NewHTTPServiceWithOptions(app.Config.Get("SAMPLE_SERVICE"), app.Logger, nil)

	service := services.New(sampleSvc)
	handler := handlers.New(service)

	app.GET("/user/{name}", handler.Get)

	app.Start()
}
