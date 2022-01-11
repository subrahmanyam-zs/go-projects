package main

import (
	"sync"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/service"
)

func main() {
	app := gofr.New()

	url := app.Config.Get("SAMPLE_API_URL")

	app.GET("/trace", func(ctx *gofr.Context) (interface{}, error) {
		span := ctx.Trace("some-sample-work")
		<-time.After(time.Millisecond * 1) // Waiting for 1ms to simulate workload
		span.End()

		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			svc := service.NewHTTPServiceWithOptions(url, ctx.Logger, nil)
			_, _ = svc.Get(ctx, "hello", nil)
			wg.Done()
		}()

		// Ping redis 2 times concurrently and wait.
		count := 2
		wg.Add(count)
		for i := 0; i < count; i++ {
			go func() {
				ctx.Redis.Ping(ctx)
				wg.Done()
			}()
		}
		wg.Wait()

		return "ok", nil
	})

	app.Start()
}
