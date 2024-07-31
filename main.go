package main

import (
	"alfa/api"
	"alfa/cache"
	"alfa/data"
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

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

	api.GetAllBreedsEndpoint(e)
	api.GetBreedImages(e)

	e.Logger.Fatal(e.Start(":8000"))
}
