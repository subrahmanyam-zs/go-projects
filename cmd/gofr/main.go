package main

import (
	addroute "developer.zopsmart.com/go/gofr/cmd/gofr/addRoute"
	"developer.zopsmart.com/go/gofr/cmd/gofr/entity"
	"developer.zopsmart.com/go/gofr/cmd/gofr/initialize"
	"developer.zopsmart.com/go/gofr/cmd/gofr/migration/handler"
	"developer.zopsmart.com/go/gofr/cmd/gofr/test"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	k := gofr.NewCMD()

	k.GET("migrate create", handler.CreateMigration)
	k.GET("migrate", handler.Migrate)
	k.GET("init", initialize.Init)
	k.GET("entity", entity.AddEntity)
	k.GET("add", addroute.AddRoute)
	k.GET("help", helpHandler)
	k.GET("test", test.GenerateIntegrationTest)
	k.Start()
}

func helpHandler(c *gofr.Context) (interface{}, error) {
	return `Available Commands
init
entity
add
test
migrate
migrate create

Run gofr <command_name> -h for help of the command`, nil
}
