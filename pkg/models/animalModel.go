package models

import (
	"gorm.io/gorm"
)

type Animal struct {
	gorm.Model
	Name        string
	Age         float32
	Type        string
	Description string
	Gender      string
	Vaccinated  bool
	Sterilized  bool
	Image       Image
	Photos      []Photo
}

// Users       []User `gorm:"many2many:seen_animals;"`

type AnimalJSON struct {
	ID          uint     `json:"id"`
	Name        string   `json:"name" binding:"required,alphanum,min=1,max=30"`
	Age         float32  `json:"age" binding:"required,numeric,min=0,max=30"`
	Type        string   `json:"type" binding:"required,min=1,max=30"`
	Description string   `json:"description" binding:"required,max=400"`
	Gender      string   `json:"gender" binding:"required,uppercase,contains,min=1,max=30"`
	Vaccinated  bool     `json:"vaccinated"  binding:"boolean"`
	Sterilized  bool     `json:"sterilized"  binding:"boolean"`
	Image       string   `json:"image" `
	Photos      []string `json:"photos"`
}

func ToAnimalJSON(a Animal) AnimalJSON {
	result := AnimalJSON{
		ID:          a.ID,
		Name:        a.Name,
		Age:         a.Age,
		Type:        a.Type,
		Description: a.Description,
		Gender:      a.Gender,
		Vaccinated:  a.Vaccinated,
		Sterilized:  a.Sterilized,
		Image:       a.Image.URL,
		Photos:      PhotosToArray(a.Photos),
	}
	return result
}

func FromAnimalJSON(a *AnimalJSON) *Animal {
	result := Animal{
		Name:        a.Name,
		Age:         a.Age,
		Type:        a.Type,
		Description: a.Description,
		Gender:      a.Gender,
		Vaccinated:  a.Vaccinated,
		Sterilized:  a.Sterilized,
	}
	return &result
}

func ToAnimalJSONArray(data []Animal) []AnimalJSON {
	animals := []AnimalJSON{}
	for _, a := range data {
		animals = append(animals, ToAnimalJSON(a))
	}
	return animals
}
