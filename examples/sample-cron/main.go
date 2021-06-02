package main

import (
	"log"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

var n = 0

func main() {
	c := gofr.NewCron()

	// runs every minute
	c.AddJob("* * * * *", count)

	// setting maximum duration of this program
	time.Sleep(3 * time.Minute)
}

func count() {
	n++
	log.Println("Count: ", n)
}
