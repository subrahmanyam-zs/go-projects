package handler

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/template"
)

func FileHandler(ctx *gofr.Context) (interface{}, error) {
	return template.Template{Directory: ctx.TemplateDir, File: "gofr.png", Data: nil, Type: template.FILE}, nil
}
