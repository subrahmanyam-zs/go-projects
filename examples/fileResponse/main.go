package main

import (
	"os"

	"developer.zopsmart.com/go/gofr/examples/fileResponse/handler"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	// Create the application object
	k := gofr.New()
	rootPath, _ := os.Getwd()

	// overriding default template location.
	k.TemplateDir = rootPath + "/static"

	k.GET("/file", handler.FileHandler)

	k.Start()
}
