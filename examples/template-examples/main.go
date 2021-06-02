package main

import (
	"os"

	"developer.zopsmart.com/go/gofr/examples/template-examples/handler"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	// Create the application object
	k := gofr.New()
	rootPath, _ := os.Getwd()

	// overriding default template location.
	k.TemplateDir = rootPath + "/templates"

	k.GET("/test", handler.Template)

	k.GET("/image", handler.Image)

	k.Start()
}
