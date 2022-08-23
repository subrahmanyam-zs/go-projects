package template

import (
	"bytes"
	"html/template"
	"mime"
	"os"
	"path/filepath"

	"developer.zopsmart.com/go/gofr/pkg/errors"
)

type fileType int

const (
	HTML fileType = iota
	TEXT
	CSV
	FILE
)

type File struct {
	Content     []byte
	ContentType string
}

type Template struct {
	Directory string
	File      string
	Data      interface{}
	Type      fileType
}

func (t *Template) Render() ([]byte, error) {
	defaultLocation := t.Directory
	// if the temp location is not specified
	// the default location is taken from root of the project
	if defaultLocation == "" {
		rootLocation, _ := os.Getwd()
		defaultLocation = rootLocation + "/static"
	}

	templ, err := template.New(t.File).ParseFiles(defaultLocation + "/" + t.File)
	if err != nil {
		return nil, errors.FileNotFound{Path: t.Directory, FileName: t.File}
	}

	if t.Data != nil {
		var tpl bytes.Buffer

		err = templ.Execute(&tpl, t.Data)
		if err != nil {
			return nil, err
		}

		return tpl.Bytes(), nil
	}

	path := defaultLocation + "/" + t.File

	b, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (t *Template) ContentType() string {
	switch t.Type {
	case HTML:
		return "text/html"
	case CSV:
		return "text/csv"
	case TEXT:
		return "text/plain"
	case FILE:
		extn := filepath.Ext(t.File)
		if extn == ".json" {
			return "application/json"
		}

		return mime.TypeByExtension(extn)
	default:
		return "text/plain"
	}
}
