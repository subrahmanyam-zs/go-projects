package main

import (
	"developer.zopsmart.com/go/gofr/examples/using-file/handler"
	"developer.zopsmart.com/go/gofr/pkg/file"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	app := gofr.NewCMD()

	fileAbstracter, err := file.NewWithConfig(app.Config, "test.txt", "rw")
	if err != nil {
		app.Logger.Error("Unable to initialize", err)
		return
	}

	h := handler.New(fileAbstracter)

	app.GET("read", h.Read)
	app.GET("write", h.Write)
	app.GET("list", h.List)

	app.Start()
}
