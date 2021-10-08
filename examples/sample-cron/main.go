package main

import (
	"log"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

//nolint:gochecknoglobals // used in main_test.go
var n = 0

const minute = 3

func main() {
	app := gofr.New()

	c := gofr.NewCron()

	// runs every minute
	err := c.AddJob("* * * * *", count)
	if err != nil {
		app.Logger.Error(err)
		return
	}

	// setting maximum duration of this program
	time.Sleep(minute * time.Minute)
}

func count() {
	n++
	log.Println("Count: ", n)
}
