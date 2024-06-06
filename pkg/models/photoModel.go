package models

import (
	"gorm.io/gorm"
)

type Photo struct {
	gorm.Model
	AnimalID uint
	ImageURL string
	Key      string // s3 upload id
}

func PhotosToArray(photos []Photo) (urls []string) {
	for _, p := range photos {
		urls = append(urls, p.ImageURL)
	}
	return urls
}
