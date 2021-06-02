package gofr

import (
	"io/fs"
	"os"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr/template"
	"developer.zopsmart.com/go/gofr/web"
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
