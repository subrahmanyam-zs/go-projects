package main

import (
	"os"

	"developer.zopsmart.com/go/gofr/examples/template-examples/handler"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	// Create the application object
	app := gofr.New()
	rootPath, _ := os.Getwd()

	// overriding default template location.
	app.TemplateDir = rootPath + "/templates"

	app.GET("/test", handler.Template)

	app.GET("/image", handler.Image)

	app.Start()
}
