package web

import (
	"embed"

	"developer.zopsmart.com/go/gofr/pkg/gofr/template"
)

//go:embed swagger/*
var fs embed.FS

func GetSwaggerFile(fileName string) (data []byte, contentType string, err error) {
	t := template.Template{}
	if fileName == "" {
		t.File = "index.html"
		t.Type = template.HTML
	} else {
		t.File = fileName
		t.Type = template.FILE
	}

	data, err = fs.ReadFile("swagger/" + t.File)
	if err != nil {
		return
	}

	contentType = t.ContentType()

	return
}
