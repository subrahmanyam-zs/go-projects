package handler

import (
	"fmt"
	"sync"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/service"
)

// HelloWorld is a handler function of type gofr.Handler, it responds with a message
func HelloWorld(c *gofr.Context) (interface{}, error) {
	return "Hello World!", nil
}

// HelloName is a handler function of type gofr.Handler, it responds with a message and uses query params
func HelloName(c *gofr.Context) (interface{}, error) {
	return fmt.Sprintf("Hello %s", c.Param("name")), nil
}

// ErrorHandler always returns an error
func ErrorHandler(c *gofr.Context) (interface{}, error) {
	return nil, &errors.Response{
		StatusCode: 500,
		Code:       "UNKNOWN_ERROR",
		Reason:     "unknown error occurred",
	}
}

type resp struct {
	Name    string
	Company string
}

// JSONHandler is a handler function of type gofr.Handler, it responds with a JSON message
func JSONHandler(c *gofr.Context) (interface{}, error) {
	r := resp{
		Name:    "Vikash",
		Company: "ZopSmart",
	}

	return r, nil
}

type handler struct {
	serviceUrl string
}

func New(url string) handler {
	return handler{serviceUrl: url}
}

func (h handler) Trace(c *gofr.Context) (interface{}, error) {
	span2 := c.Trace("some-sample-work")
	<-time.After(time.Millisecond * 1) // Waiting for 1ms to simulate workload
	span2.End()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		svc := service.NewHTTPServiceWithOptions(h.serviceUrl, c.Logger, nil)
		_, _ = svc.Get(c, "config/key", nil)
		wg.Done()
	}()

	// Ping redis 5 times concurrently and wait.
	count := 10
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			c.Redis.Ping(c)
			wg.Done()
		}()
	}
	wg.Wait()

	return "ok", nil
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
