package repository

import (
	"alfa/cache"
	"alfa/utils"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

type dogs struct{}
type resultChan struct {
	images []string
}

func CacheResponse(key string, response utils.Response) error {
	rdb := cache.GetRedisClient()
	ctx := cache.GetRedisCtx()

	data, err := json.Marshal(response)
	if err != nil {
		return err
	}

	err = rdb.Set(ctx, key, data, 30*time.Second).Err()
	if err != nil {
		return err
	}

	return nil
}

func CacheGenericResponse(key string, response utils.GenericResponse) error {
	rdb := cache.GetRedisClient()
	ctx := cache.GetRedisCtx()

	data, err := json.Marshal(response)
	if err != nil {
		return err
	}

	err = rdb.Set(ctx, key, data, 30*time.Second).Err()
	if err != nil {
		return err
	}

	return nil
}

func GetCachedResponse(key string) (*utils.Response, error) {
	rdb := cache.GetRedisClient()
	ctx := cache.GetRedisCtx()

	cdata, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var response utils.Response
	if err := json.Unmarshal([]byte(cdata), &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func GetCachedImage(key string) (*utils.GenericResponse, error) {
	rdb := cache.GetRedisClient()
	ctx := cache.GetRedisCtx()

	cdata, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var response utils.GenericResponse
	if err := json.Unmarshal([]byte(cdata), &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func GetAllBreeds(c echo.Context) error {
	var url = os.Getenv("URL")
	var breedsResponse utils.Response
	breedTypes := make(map[string][]string)

	ListAllBreedsConfig := utils.RequestsConfig{
		URL:     url + "/breeds/list/all",
		Timeout: 5 * time.Second,
	}

	cdata, err := GetCachedResponse("allBreeds")
	if err != nil {
		log.Printf("Error getting allBreeds: %v", err)
	}

	if cdata != nil {
		log.Print("Using cached data")
		return c.JSON(http.StatusOK, cdata)
	}

	breedTypesDB := GetBreedTypes()
	if breedTypesDB != nil {
		breedsResponse.Message = breedTypesDB
		breedsResponse.Status = "success"

		log.Printf("Data found on breeds")
	} else {
		log.Printf("No data found on breeds")
		listAllBreedRequest := utils.NewRequests(ListAllBreedsConfig)

		res, err := listAllBreedRequest.Get()
		if err != nil {
			log.Panicf("GET request failed: %v", err)
			// return c.JSON(http.StatusInternalServerError, err.Error())
			return c.JSON(http.StatusOK, []string{})
		}

		if err := json.Unmarshal(res, &breedsResponse); err != nil {
			// return c.JSON(http.StatusInternalServerError, "Failed to unmarshal JSON data")
			return c.JSON(http.StatusOK, []string{})
		}
		err_pop := PopulateDog(breedsResponse.Message)
		if err_pop != nil {
			log.Printf("Error populating dog table: %v", err_pop)
		}
		log.Print("Dog table populated")
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

	go fetchImgChan(c, "terrier", results)
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
	errc := CacheResponse("allBreeds", modifiedResponse)
	if errc != nil {
		log.Printf("Error caching dog_all: %v", errc)
	}
	return c.JSON(http.StatusOK, modifiedResponse)
}

func fetchImgChan(c echo.Context, breedType string, results chan<- resultChan) {
	var url = os.Getenv("URL")
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
	results <- resultChan{breedImgResponse.Message}
	close(results)
	return
}

func GetBreedImages(c echo.Context) error {
	breed := c.Param("breed")
	var breedImgResponse utils.GenericResponse
	var filteredImages []string
	var url = os.Getenv("URL")
	key := "images" + breed

	cdata, err := GetCachedImage(key)
	if err != nil {
		log.Printf("Error getting cached images: %v", err)
	}

	if cdata != nil {
		log.Print("Using cached images")
		return c.JSON(http.StatusOK, cdata)
	}

	breedImagesConfig := utils.RequestsConfig{
		URL:     url + "/breed/" + breed + "/images",
		Timeout: 2 * time.Second,
	}

	images := GetBreedImagesList(breed)
	if images != nil {
		log.Printf("Data found on %s images", breed)
		breedImgResponse.Message = images
		breedImgResponse.Status = "success"
	} else {
		log.Printf("No data found on %s images", breed)
		breedImagesRequest := utils.NewRequests(breedImagesConfig)
		res, err := breedImagesRequest.Get()
		if err != nil {
			log.Printf("GET request failed: %v", err)
			// return c.JSON(http.StatusInternalServerError, err.Error())
			return c.JSON(http.StatusOK, []string{})
		}

		if err := json.Unmarshal(res, &breedImgResponse); err != nil {
			log.Printf("Error %v", err)
			return c.JSON(http.StatusOK, []string{})
		}

	}

	if strings.ToLower(breed) == "shiba" {
		for _, img := range breedImgResponse.Message {
			if isOddNumberInFilename(img) {
				filteredImages = append(filteredImages, img)
			}
		}
	} else {
		filteredImages = breedImgResponse.Message
	}

	if images == nil {
		pop_img := PopulateDogImage(filteredImages, breed)
		if pop_img != nil {
			log.Printf("Breed %s Images inserted", breed)
		}
	}

	breedImgResponse.Message = filteredImages
	errc := CacheGenericResponse(key, breedImgResponse)
	if errc != nil {
		log.Printf("Error caching images: %v", errc)
	}
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
