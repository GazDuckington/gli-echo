package main

import (
	"alfa/api"
	"alfa/cache"
	"alfa/data"
	"log"

	_ "alfa/docs"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/echo-swagger"
)

// @title API Docs
// @version 1.0

// @host localhost:8000
func main() {
	e := echo.New()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println("Environment loaded")

	err_r := cache.Initialize()
	if err_r != nil {
		log.Fatalf("Error init redis cache: %v", err_r)
	}
	log.Println("Redis cache loaded")

	data.Initialize()

	e.Use(middleware.Logger())

	e.GET("/swagger-ui/*", echoSwagger.WrapHandler)

	api.GetAllBreedsEndpoint(e)
	api.GetBreedImages(e)

	e.Logger.Fatal(e.Start(":8000"))
}
