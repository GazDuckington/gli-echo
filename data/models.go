package data

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Dog struct {
	ID        string     `gorm:"type:TEXT;primaryKey"`
	Breed     string     `gorm:"size:255"`
	Subbreeds JSONString `gorm:"type:TEXT"`
}

type DogImage struct {
	ID     string     `gorm:"type:TEXT;primaryKey"`
	Breed  string     `gorm:"size:255"`
	Images JSONString `gorm:"type:TEXT"`
}

func (d *Dog) BeforeCreate(tx *gorm.DB) (err error) {
	d.ID = uuid.New().String()
	return
}

func (d *DogImage) BeforeCreate(tx *gorm.DB) (err error) {
	d.ID = uuid.New().String()
	return
}

type JSONString []string

func (j *JSONString) Scan(value interface{}) error {
	if value == nil {
		*j = JSONString{}
		return nil
	}

	switch v := value.(type) {
	case string:
		if err := json.Unmarshal([]byte(v), j); err != nil {
			return fmt.Errorf("failed to unmarshal JSON string: %v", err)
		}
	case []byte:
		if err := json.Unmarshal(v, j); err != nil {
			return fmt.Errorf("failed to unmarshal JSON byte slice: %v", err)
		}
	default:
		return fmt.Errorf("expected string or []byte but got %T", value)
	}
	return nil
}

func (j JSONString) Value() (driver.Value, error) {
	if j == nil {
		return "[]", nil
	}
	data, err := json.Marshal(j)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON data: %v", err)
	}
	return string(data), nil
}
