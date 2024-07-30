package repository

import (
	"alfa/utils"
	"encoding/json"
	"fmt"
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
	type resultChan struct {
		images []string
	}
	fetchImgChan := func(breedType string, results chan<- resultChan) {
		req := c.Request().Clone(c.Request().Context())
		rec := c.Response().Writer
		echoCtx := c.Echo().NewContext(req, rec)
		echoCtx.SetParamNames("breed")
		echoCtx.SetParamValues(breedType)

		breedurl := echoCtx.Param("breed")
		breedImagesConfig := utils.RequestsConfig{
			URL:     url + "/breed/" + breedurl + "/images",
			Timeout: 2 * time.Second,
		}

		breedImagesRequest := utils.NewRequests(breedImagesConfig)
		res, err := breedImagesRequest.Get()
		if err != nil {
			log.Panicf("GET request failed: %v", err)
			results <- resultChan{nil}
			return
		}

		type restype struct {
			Message []string `json:"message"`
			Status  string   `json:"status"`
		}
		var breedImgResponse restype
		if err := json.Unmarshal(res, &breedImgResponse); err != nil {
			log.Printf("Error %v", err)
			results <- resultChan{nil}
			return
		}
		fmt.Printf("AAAA: %s", breedImgResponse)
		results <- resultChan{breedImgResponse.Message}
		close(results)
	}

	var terrierBreeds []string
	results := make(chan resultChan)

	for breedType, breeds := range breedsResponse.Message {
		if len(breeds) > 0 {
			switch strings.ToLower(breedType) {
			case "sheepdog":
				for _, breed := range breeds {
					breedTypes[breedType+"-"+breed] = []string{}
				}
				delete(breedTypes, "sheepdog")
			case "terrier":
				terrierBreeds = breeds
				delete(breedTypes, "terrier")
			default:
				breedTypes[breedType] = breeds
			}
		}
	}

	go fetchImgChan("terrier", results)
	for result := range results {
		for _, t := range terrierBreeds {
			key := "terrier" + "-" + t
			breedTypes[key] = findSubBreed(result.images, key)
		}
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
			if isOddNumberInFilename(img) {
				filteredImages = append(filteredImages, img)
			}
		}
	} else {
		filteredImages = breedImgResponse.Message
	}

	breedImgResponse.Message = filteredImages

	return c.JSON(http.StatusOK, breedImgResponse)
}

func isOddNumberInFilename(filename string) bool {
	re := regexp.MustCompile(`(\d+)\.jpg$`)
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

func findSubBreed(filenames []string, subBreed string) []string {
	pattern := ".*" + regexp.QuoteMeta(subBreed) + ".*"
	regex, err := regexp.Compile(pattern)
	if err != nil {
		log.Printf("Failed to compile regex: %v", err)
		return nil
	}

	var matchedFilenames []string
	for _, filename := range filenames {
		if regex.MatchString(filename) {
			matchedFilenames = append(matchedFilenames, filename)
		}
	}

	return matchedFilenames
}
