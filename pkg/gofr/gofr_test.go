package gofr

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zopsmart/gofr/pkg/gofr/config"
	"github.com/zopsmart/gofr/pkg/gofr/request"
	"github.com/zopsmart/gofr/pkg/gofr/types"
	"github.com/zopsmart/gofr/pkg/log"
)

const helloWorld = "Hello World!"

func TestGofr_ServeHTTP_TextResponse(t *testing.T) {
	testCases := []struct {
		// Given
		method string
		target string
		// Expectations
		response  string
		headerKey string
		headerVal string
	}{
		{"GET", "/hello", "Hello World!", "content-type", "text/plain"},               // Example 1
		{"PUT", "/hello", "Hello World!", "content-type", "text/plain"},               // Example 1
		{"POST", "/hello", "Hello World!", "content-type", "text/plain"},              // Example 1
		{"GET", "/params?name=Vikash", "Hello Vikash!", "content-type", "text/plain"}, // Example 2 with query parameters
	}

	k := New()
	// Added contextInjector middleware
	k.Server.Router.Use(k.Server.contextInjector)
	// Example 1 Handler
	k.GET("/hello", func(c *Context) (interface{}, error) {
		return helloWorld, nil
	})

	k.PUT("/hello", func(c *Context) (interface{}, error) {
		return helloWorld, nil
	})

	k.POST("/hello", func(c *Context) (interface{}, error) {
		return helloWorld, nil
	})

	// Example 2 Handler
	k.GET("/params", func(c *Context) (interface{}, error) {
		return fmt.Sprintf("Hello %s!", c.Param("name")), nil
	})

	for _, tc := range testCases {
		w := httptest.NewRecorder()
		r, _ := request.NewMock(tc.method, tc.target, nil)

		r.Header.Set("content-type", "text/plain")

		k.Server.Router.ServeHTTP(w, r)

		expectedResp := fmt.Sprintf("%v", &types.Response{Data: tc.response})

		if resp := w.Body.String(); resp != expectedResp {
			t.Errorf("Unexpected response for %s %s. \t expected: %s \t got: %s", tc.method, tc.target, expectedResp, resp)
		}

		if ctype := w.Header().Get(tc.headerKey); ctype != tc.headerVal {
			t.Errorf("Header mismatch for %s %s. \t expected: %s \t got: %s", tc.method, tc.target, tc.headerVal, ctype)
		}
	}
}

func TestGofr_StartPanic(t *testing.T) {
	k := New()

	http.DefaultServeMux = new(http.ServeMux)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				t.Errorf("Start funcs panics on function call")
			}
		}()
		k.Start()
	}()
	<-time.After(1 * time.Second)
}

func TestGofr_Start(t *testing.T) {
	// only http server should run therefore wrong config location given
	c := config.NewGoDotEnvProvider(log.NewMockLogger(os.Stderr), "../configserror")
	k := NewWithConfig(c)
	k.Server.UseMiddleware(sampleMW1)

	http.DefaultServeMux = new(http.ServeMux)

	go k.Start()
	time.Sleep(3 * time.Second)

	var returned = make(chan bool)

	go func() {
		http.DefaultServeMux = new(http.ServeMux)

		k1 := NewWithConfig(c)
		k1.Start()
		returned <- true
	}()
	time.Sleep(time.Second * 3)

	if !<-returned {
		t.Errorf("Was able to start server on port while server was already running")
	}
}

func TestGofr_EnableSwaggerUI(t *testing.T) {
	k := New()
	// Added contextInjector middleware
	k.Server.Router.Use(k.Server.contextInjector)

	k.EnableSwaggerUI()

	w := httptest.NewRecorder()
	r, _ := request.NewMock("GET", "/swagger", nil)

	k.Server.Router.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Errorf("Expected 200 but got: %v", w.Code)
	}
}

func TestGofrUseMiddleware(t *testing.T) {
	k := New()
	mws := []Middleware{
		sampleMW1,
		sampleMW2,
	}

	k.Server.UseMiddleware(mws...)

	if len(k.Server.mws) != 2 || !reflect.DeepEqual(k.Server.mws, mws) {
		t.Errorf("FAILED, Expected: %v, Got: %v", mws, k.Server.mws)
	}
}

func sampleMW1(h http.Handler) http.Handler {
	return h
}

func sampleMW2(h http.Handler) http.Handler {
	return h
}

func TestGofr_Config(t *testing.T) { // check config is properly set or not?
	logger := log.NewMockLogger(io.Discard)
	c := config.NewGoDotEnvProvider(logger, "../../config")
	expected := c.Get("APP_NAME")

	k := New()
	val := k.Config.Get("APP_NAME")

	if !reflect.DeepEqual(expected, val) {
		t.Errorf("FAILED, Expected: %v, Got: %v", expected, val)
	}
}

func TestGofr_Patch(t *testing.T) {
	testCases := []struct {
		// Given
		target string
		// Expectations
		expectedCode int
	}{
		{"/patch", 200},
		{"/", 404},
		{"/error", 500},
	}

	// Create a server with PATCH routes
	k := New()
	// Added contextInjector middleware
	k.Server.Router.Use(k.Server.contextInjector)

	k.Server.ValidateHeaders = false

	k.PATCH("/patch", func(c *Context) (interface{}, error) {
		return "success", nil
	})

	k.PATCH("/error", func(c *Context) (interface{}, error) {
		return nil, errors.New("sample")
	})

	for _, tc := range testCases {
		rr := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPatch, tc.target, nil)

		k.Server.Router.ServeHTTP(rr, r)

		assert.Equal(t, rr.Code, tc.expectedCode)
	}
}
