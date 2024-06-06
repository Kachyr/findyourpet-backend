package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type LikedAnimal struct {
	UserID    uuid.UUID `gorm:"primaryKey"`
	AnimalID  uint      `gorm:"primaryKey"`
	Timestamp time.Time
}

func (s *LikedAnimal) BeforeCreate(tx *gorm.DB) error {
	s.Timestamp = time.Now()
	log.Info().Any("TEST BeforeCreate", s.Timestamp).Send()
	return nil
}

func (s *LikedAnimal) BeforeUpdate(tx *gorm.DB) error {
	s.Timestamp = time.Now()
	log.Info().Any("TEST BeforeUpdate", s.Timestamp).Send()
	return nil
}
