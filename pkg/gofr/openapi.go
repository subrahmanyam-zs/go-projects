package gofr

import (
	"io/fs"
	"os"

	"github.com/zopsmart/gofr/pkg/errors"
	"github.com/zopsmart/gofr/pkg/gofr/template"
	"github.com/zopsmart/gofr/web"
)

// OpenAPIHandler serves the openapi.json file present either in the root directory or in root/api directory
func OpenAPIHandler(c *Context) (interface{}, error) {
	rootDir, _ := os.Getwd()
	fileDir := rootDir + "/" + "api"

	return template.Template{Directory: fileDir, File: "openapi.json", Data: nil, Type: template.FILE}, nil
}

func SwaggerUIHandler(c *Context) (interface{}, error) {
	fileName := c.PathParam("name")

	data, contentType, err := web.GetSwaggerFile(fileName)
	if err != nil {
		switch err.(type) {
		case *fs.PathError:
			return nil, errors.FileNotFound{
				FileName: fileName,
			}
		default:
			return nil, err
		}
	}

	return template.File{Content: data, ContentType: contentType}, nil
}
