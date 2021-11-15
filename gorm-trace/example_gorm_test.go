package otelgorm

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func ExampleNewPlugin() {
	// Initialize db connection
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Initialize otel plugin with options
	plugin := NewPlugin(
	// include any options here
	)
	err = db.Use(plugin)
	if err != nil {
		panic("failed configuring plugin")
	}

	// Migrate the schema
	err = db.AutoMigrate(&Product{})
	if err != nil {
		panic(err.Error())
	}

	// Create
	db.Create(&Product{Code: "D42", Price: 100})

	// Read
	var product Product
	db.First(&product, 1)                 // find product with integer primary key
	db.First(&product, "code = ?", "D42") // find product with code D42
}
