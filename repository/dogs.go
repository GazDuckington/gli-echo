package repository

import (
	"alfa/utils"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type dogs struct{}

func GetAllBreeds(c echo.Context) error {
	var url = os.Getenv("URL")
	ListAllBreedsConfig := utils.RequestsConfig{
		URL:     url + "/breeds/list/all",
		Timeout: 5 * time.Second,
	}

	listAllBreedRequest := utils.NewRequests(ListAllBreedsConfig)

	res, err := listAllBreedRequest.Get()
	if err != nil {
		log.Panicf("GET request failed: %v", err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	var breedsResponse utils.Response
	if err := json.Unmarshal(res, &breedsResponse); err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to unmarshal JSON data")
	}

	breedTypes := make(map[string][]string)
	for breedType, breeds := range breedsResponse.Message {
		if len(breeds) > 0 {
			if strings.ToLower(breedType) == "sheepdog" {
				for _, breed := range breeds {
					breedTypes[breedType+"-"+breed] = []string{}
				}
			}
			breedTypes[breedType] = breeds
		}
		breedTypes[breedType] = breeds

	}

	modifiedResponse := utils.Response{
		Message: breedTypes,
		Status:  breedsResponse.Status,
	}
	return c.JSON(http.StatusOK, modifiedResponse)
}

func GetBreedImages(c echo.Context) error {
	breed := c.Param("breed")

	var url = os.Getenv("URL")
	breedImagesConfig := utils.RequestsConfig{
		URL:     url + "/breed/" + breed + "/images",
		Timeout: 2 * time.Second,
	}

	breedImagesRequest := utils.NewRequests(breedImagesConfig)
	res, err := breedImagesRequest.Get()
	if err != nil {
		log.Panicf("GET request failed: %v", err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	var breedImgResponse utils.GenericResponse
	if err := json.Unmarshal(res, &breedImgResponse); err != nil {
		log.Panicf("Error %v", err)
		return c.JSON(http.StatusInternalServerError, "Failed to unmarshal JSON data")
	}

	var filteredImages []string
	if strings.ToLower(breed) == "shiba" {
		for _, img := range breedImgResponse.Message {
			if strings.ToLower(breed) == "shiba" {
				if isOddNumberInFilename(img) {
					filteredImages = append(filteredImages, img)
				}
			}
		}
	} else {
		filteredImages = breedImgResponse.Message
	}

	breedImgResponse.Message = filteredImages

	return c.JSON(http.StatusOK, breedImgResponse)
}
func isOddNumberInFilename(filename string) bool {
	re := regexp.MustCompile(`(\d+)\.jpg$`) // Regex to extract number before .jpg
	matches := re.FindStringSubmatch(filename)
	if len(matches) > 1 {
		numberStr := matches[1]
		number, err := strconv.Atoi(numberStr)
		if err == nil && number%2 != 0 {
			return true
		}
	}
	return false
}
