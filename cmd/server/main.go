package main

import (
	"net/http"

	"github.com/Mazzael/go-api/configs"
	"github.com/Mazzael/go-api/internal/entity"
	"github.com/Mazzael/go-api/internal/infra/database"
	"github.com/Mazzael/go-api/internal/infra/webserver/handlers"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	_, err := configs.LoadConfig(".")
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database: " + err.Error())
	}

	db.AutoMigrate(&entity.Product{}, &entity.User{})

	gormProductRepository := database.NewProduct(db)

	productHandler := handlers.NewProductHandler(gormProductRepository)

	http.HandleFunc("/products", productHandler.CreateProduct)
	http.ListenAndServe(":8080", nil)
}
