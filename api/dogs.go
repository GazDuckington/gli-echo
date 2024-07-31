package api

import (
	"alfa/repository"
	"github.com/labstack/echo/v4"
)

// @Summary Get dog breeds
// @Description Get all breeds
// @Tags breeds
// @Accept json
// @Produce json
// @Success 200 {array} string
// @Router /breeds/ [get]
func GetAllBreedsEndpoint(e *echo.Echo) {
	e.GET("/", repository.GetAllBreeds)
	e.GET("/breeds/", repository.GetAllBreeds)
}

// @Summary Get dog images
// @Description Get images of a specific breed
// @Tags breed-images
// @Accept json
// @Produce json
// @Success 200 {array} string
// @Param breed path string true "Breed"
// @Router /breeds/{breed}/images [get]
func GetBreedImages(e *echo.Echo) {
	e.GET("/breeds/:breed/images", repository.GetBreedImages)
	e.GET("/:breed/images", repository.GetBreedImages)
}
