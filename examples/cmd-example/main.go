package main

import (
	"errors"
	"fmt"

	"github.com/zopsmart/gofr/pkg/gofr"
	"github.com/zopsmart/gofr/pkg/gofr/template"
)

func main() {
	k := gofr.NewCMD()

	k.GET("hello", func(c *gofr.Context) (i interface{}, err error) {
		name := c.PathParam("name")
		if name == "" {
			return fmt.Sprint("Hello!"), nil
		}
		return fmt.Sprintf("Hello %s!", name), nil
	})

	k.GET("error", func(c *gofr.Context) (i interface{}, err error) {
		return nil, errors.New("some error occurred")
	})

	k.GET("bind", func(c *gofr.Context) (i interface{}, err error) {
		var a struct {
			Name   string
			IsGood bool
		}

		_ = c.Bind(&a)

		return fmt.Sprintf("Name: %s Good: %v", a.Name, a.IsGood), nil
	})

	k.GET("temp", func(c *gofr.Context) (i interface{}, err error) {
		filename := c.PathParam("filename")
		return template.Template{
			Directory: "templates",
			File:      filename,
			Data:      nil,
			Type:      template.FILE,
		}, nil
	})

	k.GET("file", func(c *gofr.Context) (i interface{}, err error) {
		return template.File{
			Content:     []byte("Hello"),
			ContentType: "text/plain",
		}, nil
	})

	k.Start()
}
