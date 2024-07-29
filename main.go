package main

import (
	"alfa/api"
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println("Environment loaded")

	api.GetAllBreedsEndpoint(e)
	api.GetBreedImages(e)

	e.Logger.Fatal(e.Start(":8080"))
}
