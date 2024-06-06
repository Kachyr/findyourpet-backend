package users

import (
	"github.com/Kachyr/findyourpet/findyourpet-backend/pkg/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserStoreI interface {
	Create(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	GetByID(id uuid.UUID) (*models.User, error)
	GetUserSettings(id uuid.UUID) (models.UserSettings, error)
	SetUserSettings(userID uuid.UUID, newSettings models.UserSettings) error
}

type UserStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) Create(user *models.User) error {
	result := s.db.Create(&user)

	return result.Error
}

func (s *UserStore) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	result := s.db.First(&user, "email = ?", email)
	return user, result.Error
}

func (s *UserStore) GetByID(id uuid.UUID) (*models.User, error) {
	user := &models.User{}
	result := s.db.First(&user, id)
	return user, result.Error
}

func (s *UserStore) SetUserSettings(userID uuid.UUID, newSettings models.UserSettings) error {
	var settings models.UserSettings
	if err := s.db.First(&settings, "user_id = ?", userID).Error; err != nil {
		return err
	}

	settings.Gender = newSettings.Gender
	settings.Location = newSettings.Location
	settings.MaxAge = newSettings.MaxAge
	settings.MinAge = newSettings.MinAge
	settings.Sterilized = newSettings.Sterilized
	settings.Type = newSettings.Type
	settings.Vaccinated = newSettings.Vaccinated

	return s.db.Save(&settings).Error
}

func (s *UserStore) GetUserSettings(id uuid.UUID) (models.UserSettings, error) {
	user := &models.User{}
	result := s.db.Preload("UserSettings").First(&user, id)
	return user.UserSettings, result.Error
}
