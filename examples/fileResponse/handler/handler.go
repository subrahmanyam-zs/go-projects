package handler

import (
	"github.com/zopsmart/gofr/pkg/gofr"
	"github.com/zopsmart/gofr/pkg/gofr/template"
)

func FileHandler(c *gofr.Context) (interface{}, error) {
	return template.Template{Directory: c.TemplateDir, File: "gofr.png", Data: nil, Type: template.FILE}, nil
}
