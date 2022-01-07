package main

import (
	"net/http"
	"time"

	handlers "developer.zopsmart.com/go/gofr/examples/using-http-service/handlers/user"
	services "developer.zopsmart.com/go/gofr/examples/using-http-service/services/user"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/log"
	svc "developer.zopsmart.com/go/gofr/pkg/service"
)

func main() {
	app := gofr.New()

	const numOfRetries = 3

	sampleSvc := svc.NewHTTPServiceWithOptions(app.Config.Get("SAMPLE_SERVICE"), app.Logger, &svc.Options{
		NumOfRetries: numOfRetries,
	})

	service := services.New(sampleSvc)
	handler := handlers.New(service)

	app.GET("/user/{name}", handler.Get)

	// custom retry logic to retry the service call on error or non 200 status code.
	sampleSvc.CustomRetry = func(logger log.Logger, err error, statusCode, attemptCount int) bool {
		if statusCode == http.StatusOK {
			return false
		}

		if err != nil {
			logger.Logf("got error %v, on attempt %v", err, attemptCount)
		}

		//nolint:gomnd // introducing constants for attemptCount values will reduce readability.
		switch attemptCount {
		case 1:
			time.Sleep(2 * time.Second)
		case 2:
			time.Sleep(4 * time.Second)
		case 3:
			time.Sleep(8 * time.Second)
		}

		return true
	}

	app.Start()
}
