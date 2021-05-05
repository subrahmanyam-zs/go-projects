package main

import (
	addroute "github.com/zopsmart/gofr/cmd/gofr/addRoute"
	"github.com/zopsmart/gofr/cmd/gofr/entity"
	"github.com/zopsmart/gofr/cmd/gofr/initialize"
	"github.com/zopsmart/gofr/cmd/gofr/migration/handler"
	"github.com/zopsmart/gofr/cmd/gofr/test"
	"github.com/zopsmart/gofr/pkg/gofr"
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
