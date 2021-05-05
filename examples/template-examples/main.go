package main

import (
	"os"

	"github.com/zopsmart/gofr/examples/template-examples/handler"
	"github.com/zopsmart/gofr/pkg/gofr"
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
