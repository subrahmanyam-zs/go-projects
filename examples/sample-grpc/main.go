package main

import (
	"github.com/zopsmart/gofr/examples/sample-grpc/handler/grpc"
	"github.com/zopsmart/gofr/examples/sample-grpc/handler/http"
	"github.com/zopsmart/gofr/pkg/gofr"
)

func main() {
	// this example shows an applicationZ that uses both, HTTP and GRPC
	k := gofr.New()
	k.GET("/example", http.Get)
	grpc.RegisterExampleServiceServer(k.Server.GRPC.Server(), grpc.Handler{})

	k.Start()
}
