package main

import (
	"crypto/tls"
	"net/http"

	avroHandler "developer.zopsmart.com/go/gofr/examples/universal-example/avro/handlers"
	cassandraHandler "developer.zopsmart.com/go/gofr/examples/universal-example/cassandra/handlers"
	cassandraStore "developer.zopsmart.com/go/gofr/examples/universal-example/cassandra/store/employee"
	eventHandler "developer.zopsmart.com/go/gofr/examples/universal-example/eventhub/handlers"
	svcHandler "developer.zopsmart.com/go/gofr/examples/universal-example/gofr-services/handler"
	svc "developer.zopsmart.com/go/gofr/examples/universal-example/gofr-services/service"
	pgsqlHandler "developer.zopsmart.com/go/gofr/examples/universal-example/pgsql/handler"
	pgsqlStore "developer.zopsmart.com/go/gofr/examples/universal-example/pgsql/store"
	redisHandler "developer.zopsmart.com/go/gofr/examples/universal-example/redis/handler"
	redisStore "developer.zopsmart.com/go/gofr/examples/universal-example/redis/store"
	"developer.zopsmart.com/go/gofr/pkg/datastore/pubsub/eventhub"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/service"
)

func main() {
	// Create the application object
	k := gofr.New()

	// Service urls
	urlHelloAPI := k.Config.Get("GOFR_HELLO_API")
	urlLoggingService := k.Config.Get("GOFR_LOGGING_SERVICE")

	// Skip TLS verification
	var tr = &http.Transport{
		//nolint:gosec // need this to skip TLS verification
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Gofr-logging service
	logService := service.NewHTTPServiceWithOptions(urlLoggingService, k.Logger, nil)

	logService.Client = &http.Client{Transport: tr}
	logService.Transport = tr

	logSrv := svc.New(logService)
	loggingHandler := svcHandler.New(logSrv)
	k.GET("/level", loggingHandler.Log)

	// Gofr-hello-api service
	helloService := service.NewHTTPServiceWithOptions(urlHelloAPI, k.Logger, nil)

	helloService.Client = &http.Client{Transport: tr}
	helloService.Transport = tr

	helloSrv := svc.New(helloService)
	helloSrvHandler := svcHandler.New(helloSrv)
	k.GET("/hello", helloSrvHandler.Hello)

	// Redis
	redisStr := redisStore.New()
	redisHandle := redisHandler.New(redisStr)

	k.GET("/redis/config/{key}", redisHandle.GetKey)
	k.POST("/redis/config", redisHandle.SetKey)

	// Postgres
	pgsqlStr := pgsqlStore.New()
	pgsqlHandle := pgsqlHandler.New(pgsqlStr)

	k.GET("/pgsql/employee", pgsqlHandle.Get)
	k.POST("/pgsql/employee", pgsqlHandle.Create)

	// Cassandra
	cassandraStr := cassandraStore.New()
	cassandraHandle := cassandraHandler.New(cassandraStr)

	k.GET("/cassandra/employee", cassandraHandle.Get)
	k.POST("/cassandra/employee", cassandraHandle.Create)

	k.GET("/avro/pub", avroHandler.Producer)
	k.GET("/avro/sub", avroHandler.Consumer)

	config := eventhub.Config{
		Namespace:    k.Config.Get("EVENTHUB_NAMESPACE"),
		EventhubName: k.Config.Get("EVENTHUB_NAME"),
		ClientID:     k.Config.Get("AZURE_CLIENT_ID"),
		ClientSecret: k.Config.Get("AZURE_CLIENT_SECRET"),
		TenantID:     k.Config.Get("AZURE_TENANT_ID"),
	}

	// Eventhub
	eventHub, err := eventhub.New(&config)

	if err != nil {
		k.Logger.Errorf("Azure Eventhub could not be initialized, Namespace: %v, Eventhub: %v, error: %v\n",
			config.Namespace, config.EventhubName, err)
		return
	}

	k.Logger.Infof("Azure Eventhub initialized, Namespace: %v, Eventhub: %v\n", config.Namespace, config.EventhubName)

	eventHandle := eventHandler.New(eventHub)
	k.GET("/eventhub/pub", eventHandle.Producer)
	k.GET("/eventhub/sub", eventHandle.Consumer)

	// HealthCheck for Services
	k.ServiceHealth = append(k.ServiceHealth, helloService.HealthCheck, logService.HealthCheck)

	// HealthCheck for EventHub
	k.DatabaseHealth = append(k.DatabaseHealth, eventHub.HealthCheck)

	// Start the server
	k.Start()
}
