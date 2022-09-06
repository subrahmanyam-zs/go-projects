package handler

import (
	"developer.zopsmart.com/go/gofr/examples/sample-api/entity"
	"developer.zopsmart.com/go/gofr/examples/sample-api/service"
	"fmt"
	"net/http"
	"strings"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Handler struct {
	service service.Employee
}

func New(employee service.Employee) Handler {
	return Handler{service: employee}
}

// HelloWorld is a handler function of type gofr.Handler, it responds with a message
func HelloWorld(ctx *gofr.Context) (interface{}, error) {
	return "Hello World!", nil
}

// HelloName is a handler function of type gofr.Handler, it responds with a message and uses query params
func HelloName(ctx *gofr.Context) (interface{}, error) {
	return fmt.Sprintf("Hello %s", ctx.Param("name")), nil
}

//
func Hard(ctx *gofr.Context) (interface{}, error) {
	name := ctx.Param("name")
	ctx.Logger.Debug("fetching name from request")
	ctx.Logger.Warn("warning")
	//	ctx.Logger.Fatal("fatal")
	if name == "" {
		return nil, errors.InvalidParam{Param: []string{"name"}}
	}
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

func HelloLogHandler(ctx *gofr.Context) (interface{}, error) {
	ctx.Log("key", "value")          // This is how we can add more data to framework log.
	ctx.Logger.Log("Hello Logging!") // This is how we can add a log from handlers.
	ctx.Log("key2", "value2")
	ctx.Logger.Warn("Warning 1", "Warning 2", struct {
		key1 string
		key2 int
	}{"Struct Test", 1}) // This is how you can give multiple messages

	return "Logging OK", nil
}

func (h Handler) Post(ctx *gofr.Context) (interface{}, error) {
	var emp entity.Employee

	err := ctx.Bind(&emp)
	if err != nil {
		return entity.Employee{}, errors.InvalidParam{Param: []string{"invalid"}}
	}

	res, err := h.service.Post(ctx, emp)
	if err != nil {
		return entity.Employee{}, errors.InvalidParam{Param: []string{"invalid"}}
	}

	return res, nil
}

func (h Handler) Put(ctx *gofr.Context) (interface{}, error) {
	var emp entity.Employee

	id := ctx.PathParam("id")

	err := ctx.Bind(&emp)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"invalid"}}
	}

	res, err := h.service.Put(ctx, id, emp)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"invalid"}}
	}

	return res, nil
}

func (h Handler) Delete(ctx *gofr.Context) (interface{}, error) {
	id := ctx.PathParam("id")

	res, err := h.service.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h Handler) Get(ctx *gofr.Context) (interface{}, error) {
	id := ctx.PathParam("id")

	res, err := h.service.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h Handler) GetAll(ctx *gofr.Context) (interface{}, error) {
	res, err := h.service.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}
