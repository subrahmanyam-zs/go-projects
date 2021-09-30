package main

import (
	"developer.zopsmart.com/go/gofr/examples/sample-grpc/handler/grpc"
	"developer.zopsmart.com/go/gofr/examples/sample-grpc/handler/http"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	// this example shows an applicationZ that uses both, HTTP and GRPC
	app := gofr.New()

	app.GET("/example", http.Get)

	grpcHandler := grpc.New()

	grpc.RegisterExampleServiceServer(app.Server.GRPC.Server(), grpcHandler)

	app.Start()
}
