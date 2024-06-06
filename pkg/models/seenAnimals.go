package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SeenAnimal struct {
	UserID   uuid.UUID `gorm:"primaryKey"`
	AnimalID uint      `gorm:"primaryKey"`
	Liked    bool
	SeenAt   time.Time
}

func (s *SeenAnimal) BeforeCreate(tx *gorm.DB) error {
	s.SeenAt = time.Now()

	return nil
}

func (s *SeenAnimal) BeforeUpdate(tx *gorm.DB) error {
	s.SeenAt = time.Now()

	return nil
}
