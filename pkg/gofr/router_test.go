package gofr

import (
	"fmt"
	"reflect"
	"testing"
)

func TestRouteLog(t *testing.T) {
	k := New()

	k.GET("/", func(c *Context) (interface{}, error) { return helloWorld, nil })
	k.GET("/hello-world", func(c *Context) (interface{}, error) { return helloWorld, nil })
	k.GET("/hello-world/", func(c *Context) (interface{}, error) { return helloWorld, nil })
	k.POST("/hello-world", func(c *Context) (interface{}, error) { return helloWorld, nil })
	k.POST("/hello-world/", func(c *Context) (interface{}, error) { return helloWorld, nil })
	k.POST("/hello", func(c *Context) (interface{}, error) { return helloWorld, nil })
	k.POST("/hello/", func(c *Context) (interface{}, error) { return helloWorld, nil })

	// should not be returned from logRoutes() as method is invalid
	k.Server.Router.Route("", "/failed", func(c *Context) (interface{}, error) { return helloWorld, nil })

	// should not be returned from logRoutes() as path is invalid
	k.POST("$$$$$", func(c *Context) (interface{}, error) { return helloWorld, nil })

	expected := "GET / HEAD / GET /hello-world HEAD /hello-world POST /hello-world POST /hello "

	got := fmt.Sprintf("%s"+"", k.Server.Router)

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected: %v, got: %v", expected, got)
	}
}

func TestRouter_Prefix(t *testing.T) {
	k := New()

	k.Server.Router.Prefix("/v2")
	k.GET("/hello", func(c *Context) (i interface{}, err error) {
		return "OK", nil
	})

	expected := "GET /v2/hello HEAD /v2/hello "

	got := fmt.Sprintf("%s", k.Server.Router)

	if expected != got {
		t.Errorf("FAILED, Expected: %v, Got: %v", expected, got)
	}
}
