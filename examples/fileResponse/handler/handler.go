package handler

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/template"
)

func FileHandler(c *gofr.Context) (interface{}, error) {
	return template.Template{Directory: c.TemplateDir, File: "gofr.png", Data: nil, Type: template.FILE}, nil
}
