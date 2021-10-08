package main

import (
	"os"

	"developer.zopsmart.com/go/gofr/examples/fileResponse/handler"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	// Create the application object
	app := gofr.New()
	rootPath, _ := os.Getwd()

	// overriding default template location.
	app.TemplateDir = rootPath + "/static"

	app.GET("/file", handler.FileHandler)

	app.Start()
}
