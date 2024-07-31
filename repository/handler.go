package repository

import (
	"alfa/data"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetBreedTypes() map[string][]string {
	var dogs []data.Dog

	if err := data.DB.Find(&dogs).Error; err != nil {
		return nil
	}

	breedTypes := make(map[string][]string)

	for _, dog := range dogs {
		if _, exists := breedTypes[dog.Breed]; !exists {
			breedTypes[dog.Breed] = []string{}
		}
		breedTypes[dog.Breed] = append(breedTypes[dog.Breed], dog.Subbreeds...)
	}

	if len(breedTypes) == 0 {
		return nil
	}
	if len(dogs) == 0 {
		return nil
	}
	return breedTypes
}

func GetBreedImagesList(breed string) []string {
	var dogImage data.DogImage

	if err := data.DB.Where("breed = ?", breed).First(&dogImage).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("No images found for breed %s", breed)
		} else {
			log.Printf("Error fetching images for breed %s: %v", breed, err)
		}
		return nil
	}
	if len(dogImage.Images) == 0 {
		return nil
	}
	return dogImage.Images
}

func PopulateDog(res map[string][]string) error {
	for breed, subbreeds := range res {
		dog := data.Dog{
			ID:        uuid.New().String(),
			Breed:     breed,
			Subbreeds: data.JSONString(subbreeds),
		}
		if err := data.DB.Create(&dog).Error; err != nil {
			return err
		}
	}
	return nil
}

func PopulateDogImage(res []string, breed string) error {
	dog := data.DogImage{
		ID:     uuid.New().String(),
		Breed:  breed,
		Images: data.JSONString(res),
	}
	if err := data.DB.Create(&dog).Error; err != nil {
		return err
	}

	return nil
}
