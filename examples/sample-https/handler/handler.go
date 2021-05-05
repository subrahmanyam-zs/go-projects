package handler

import (
	"fmt"

	"github.com/zopsmart/gofr/pkg/errors"
	"github.com/zopsmart/gofr/pkg/gofr"
	"github.com/zopsmart/gofr/pkg/gofr/types"
)

type Person struct {
	Username string
	Password string
}

// HelloWorld is a handler function of type gofr.Handler, it responds with a message
func HelloWorld(c *gofr.Context) (interface{}, error) {
	return "Hello World!", nil
}

func HelloName(c *gofr.Context) (interface{}, error) {
	name := c.Param("name")

	return types.Response{
		Data: fmt.Sprintf("Hello %s", name),
		Meta: map[string]interface{}{"page": 1, "offset": 0},
	}, nil
}

func PostName(c *gofr.Context) (interface{}, error) {
	var p Person

	err := c.Bind(&p)
	if err != nil {
		return nil, err
	}

	if p.Username == "alreadyExist" {
		return p, errors.EntityAlreadyExists{}
	}

	return p, nil
}

func ErrorHandler(c *gofr.Context) (interface{}, error) {
	return nil, &errors.Response{StatusCode: 404}
}

// MultipleErrorHandler returns multiple errors and
// also sets the statusCode to 400 if id is 1 else to 500
func MultipleErrorHandler(c *gofr.Context) (interface{}, error) {
	id := c.Param("id")

	var statusCode int

	if id == "1" {
		statusCode = 400
	}

	return nil, errors.MultipleErrors{
		StatusCode: statusCode,
		Errors: []error{
			errors.InvalidParam{Param: []string{"EmailAddress"}},
			errors.MissingParam{Param: []string{"Address"}},
		}}
}
