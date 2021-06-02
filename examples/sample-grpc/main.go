package main

import (
	"developer.zopsmart.com/go/gofr/examples/sample-grpc/handler/grpc"
	"developer.zopsmart.com/go/gofr/examples/sample-grpc/handler/http"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	// this example shows an applicationZ that uses both, HTTP and GRPC
	k := gofr.New()
	k.GET("/example", http.Get)
	grpc.RegisterExampleServiceServer(k.Server.GRPC.Server(), grpc.Handler{})

	k.Start()
}
