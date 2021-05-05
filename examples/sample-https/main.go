package main

import (
	"github.com/zopsmart/gofr/examples/sample-https/handler"
	"github.com/zopsmart/gofr/pkg/gofr"
)

func main() {
	k := gofr.New()

	// add a handler
	k.GET("/hello-world", handler.HelloWorld)
	k.GET("/hello", handler.HelloName)
	k.POST("/post/", handler.PostName)
	k.GET("/error", handler.ErrorHandler)
	k.GET("/multiple-errors", handler.MultipleErrorHandler)

	// set http redirect to https
	k.Server.HTTP.RedirectToHTTPS = true

	k.Start()
}
