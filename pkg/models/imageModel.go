package models

import (
	"gorm.io/gorm"
)

type Image struct {
	gorm.Model
	AnimalID uint
	URL      string
	Key      string
}
