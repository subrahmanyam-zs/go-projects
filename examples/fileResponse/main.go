package main

import (
	"os"

	"github.com/zopsmart/gofr/examples/fileResponse/handler"
	"github.com/zopsmart/gofr/pkg/gofr"
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
