// Package test provides a command line interface for running tests for a given openapi specification.
// You can run it `gofr genit -source=path/to/openapispec.yml -host=host:port`
package test

import (
	"errors"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"developer.zopsmart.com/go/gofr/cmd/gofr/helper"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func testHelp() string {
	return helper.Generate(helper.Help{
		Example:     "gofr test -host=localhost:9000 -source=/path/to/file.yml",
		Flag:        `host provide the host along with the port, source provide the path to the yml file`,
		Usage:       "test -host=<host:port> -source=</path/to/file>",
		Description: "runs integration test for a given configuration from an yml file",
	})
}

func GenerateIntegrationTest(c *gofr.Context) (interface{}, error) {
	helpFlag := c.PathParam("h")
	helpBool, _ := strconv.ParseBool(helpFlag)

	if helpBool {
		return testHelp(), nil
	}

	sourceFile := c.PathParam("source")
	if sourceFile == "" {
		return nil, errors.New("source not specified")
	}

	host := c.PathParam("host")
	if host == "" {
		return nil, errors.New("please provide host:port")
	}

	if !strings.Contains(host, "http://") {
		host = "http://" + host
	}

	swaggerLoader := openapi3.NewSwaggerLoader()
	swaggerLoader.IsExternalRefsAllowed = true

	v, err := swaggerLoader.LoadSwaggerFromFile(sourceFile)

	if err != nil {
		return nil, err
	}

	s := Swagger{openapiSwagger: v}

	err = runTests(host, s.convertIntoIntegrationTestSchema())
	if err != nil {
		return "Test Failed!", err
	}

	return "Test Passed!", nil
}
