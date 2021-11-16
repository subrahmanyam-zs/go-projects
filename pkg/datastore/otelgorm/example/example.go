package main

import (
	"context"

	"developer.zopsmart.com/go/gofr/pkg/datastore/otelgorm"
	"github.com/uptrace/opentelemetry-go-extra/otelplay"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	ctx := context.Background()

	shutdown := otelplay.ConfigureOpentelemetry(ctx)
	defer shutdown()

	db, err := gorm.Open(postgres.New(postgres.Config{DSN: "host=localhost user=postgres password=root123 dbname=customers port=2006 sslmode=disable TimeZone=Asia/Shanghai"}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if err := db.Use(otelgorm.NewPlugin()); err != nil {
		panic(err)
	}

	var num int
	if err := db.WithContext(ctx).Raw("SELECT 42").Scan(&num).Error; err != nil {
		panic(err)
	}

	otelplay.PrintTraceID(ctx)
}
