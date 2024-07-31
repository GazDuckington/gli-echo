package main

import (
	"alfa/api"
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
	data.Initialize()

	e.Use(middleware.Logger())

	api.GetAllBreedsEndpoint(e)
	api.GetBreedImages(e)

	e.Logger.Fatal(e.Start(":8000"))
}
