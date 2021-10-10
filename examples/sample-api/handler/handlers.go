package handler

import (
	"fmt"
	"net/http"
	"strings"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

// HelloWorld is a handler function of type gofr.Handler, it responds with a message
func HelloWorld(ctx *gofr.Context) (interface{}, error) {
	return "Hello World!", nil
}

// HelloName is a handler function of type gofr.Handler, it responds with a message and uses query params
func HelloName(ctx *gofr.Context) (interface{}, error) {
	return fmt.Sprintf("Hello %s", ctx.Param("name")), nil
}

// ErrorHandler always returns an error
func ErrorHandler(ctx *gofr.Context) (interface{}, error) {
	return nil, &errors.Response{
		StatusCode: http.StatusInternalServerError,
		Code:       "UNKNOWN_ERROR",
		Reason:     "unknown error occurred",
	}
}

type resp struct {
	Name    string `json:"name"`
	Company string `json:"company"`
}

// JSONHandler is a handler function of type gofr.Handler, it responds with a JSON message
func JSONHandler(ctx *gofr.Context) (interface{}, error) {
	r := resp{
		Name:    "Vikash",
		Company: "ZopSmart",
	}

	return r, nil
}

type handler struct {
	serviceUrl string
}

func HelloLogHandler(c *gofr.Context) (interface{}, error) {
	c.Log("key", "value")          // This is how we can add more data to framework log.
	c.Logger.Log("Hello Logging!") // This is how we can add a log from handlers.
	c.Log("key2", "value2")
	c.Logger.Warn("Warning 1", "Warning 2", struct {
		key1 string
		key2 int
	}{"Struct Test", 1}) // This is how you can give multiple messages

	return "Logging OK", nil
}

// UserHandler is a handler function of type gofr.Handler, it responds with a JSON message
func UserHandler(ctx *gofr.Context) (interface{}, error) {
	name := ctx.PathParam("name")

	switch strings.ToLower(name) {
	case "vikash":
		return resp{Name: "Vikash", Company: "ZopSmart"}, nil
	default:
		return nil, errors.EntityNotFound{Entity: "user", ID: name}
	}
}
