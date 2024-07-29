package api

import (
	"alfa/repository"
	"github.com/labstack/echo/v4"
)

func GetAllBreedsEndpoint(e *echo.Echo) {
	e.GET("/", repository.GetAllBreeds)
	e.GET("/breeds/", repository.GetAllBreeds)
}

func GetBreedImages(e *echo.Echo) {
	e.GET("/breeds/:breed/images", repository.GetBreedImages)
	e.GET("/:breed/images", repository.GetBreedImages)
}
