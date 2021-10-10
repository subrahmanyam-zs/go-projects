package main

import (
	"sync"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/service"
)

func main() {
	app := gofr.New()

	app.Server.HTTP.Port = 9001

	url := app.Config.GetOrDefault("SAMPLE_API_URL", "http://localhost:9000")

	app.GET("/trace", func(c *gofr.Context) (interface{}, error) {
		span2 := c.Trace("some-sample-work")
		<-time.After(time.Millisecond * 1) // Waiting for 1ms to simulate workload
		span2.End()

		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			svc := service.NewHTTPServiceWithOptions(url, c.Logger, nil)
			_, _ = svc.Get(c, "hello", nil)
			wg.Done()
		}()

		// Ping redis 2 times concurrently and wait.
		count := 2
		wg.Add(count)
		for i := 0; i < count; i++ {
			go func() {
				c.Redis.Ping(c)
				wg.Done()
			}()
		}
		wg.Wait()

		return "ok", nil
	})

	app.Start()
}
