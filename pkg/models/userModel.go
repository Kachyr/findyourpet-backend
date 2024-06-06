package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Email        string    `gorm:"unique"`
	Password     string
	UserSettings UserSettings
	SeenAnimals  []Animal `gorm:"many2many:seen_animals;"`
}

type UserSingupJSON struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=5,max=30"`
}

type UserJSON struct {
	Email        string `json:"email"`
	UserSettings UserSettingsJSON
}
