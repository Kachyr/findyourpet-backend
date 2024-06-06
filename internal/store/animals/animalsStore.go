package animals

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/Kachyr/findyourpet/findyourpet-backend/pkg/constants"
	"github.com/Kachyr/findyourpet/findyourpet-backend/pkg/models"
	"github.com/Kachyr/findyourpet/findyourpet-backend/pkg/pagination"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type AnimalStoreI interface {
	AddAnimal(animal *models.Animal) error
	AddAnimals(animals []*models.Animal) error
	GetAllAnimals(c *gin.Context) ([]models.Animal, error)
	GetById(id string) (models.Animal, error)
	GetLikedAnimals(userID uuid.UUID, c *gin.Context) ([]models.Animal, error)
	GetNotSeenAnimals(userID uuid.UUID, c *gin.Context) ([]models.Animal, error)
	MarkAsSeen(animalID uint, userID uuid.UUID, animalLiked bool) error
}

type AnimalStore struct {
	db *gorm.DB
}

func NewAnimalStore(db *gorm.DB) *AnimalStore {
	return &AnimalStore{db: db}
}

func (s *AnimalStore) AddAnimal(animal *models.Animal) error {

	return s.db.Create(&animal).Error
}
func (s *AnimalStore) AddAnimals(animals []*models.Animal) error {

	return s.db.Create(animals).Error
}

func (s *AnimalStore) GetById(id string) (models.Animal, error) {
	animal := models.Animal{}
	result := s.db.Scopes(s.addMediaPreload).First(&animal, id)
	return animal, result.Error
}

func (s *AnimalStore) GetLikedAnimals(userID uuid.UUID, c *gin.Context) ([]models.Animal, error) {
	var animals []models.Animal
	err := s.db.
		Joins("JOIN seen_animals ON animals.id = seen_animals.animal_id").
		Where("seen_animals.user_id = ? AND seen_animals.liked = ?", userID, true).
		Scopes(s.addMediaPreload, pagination.Paginate(c)).
		Find(&animals).Error
	if err != nil {
		return nil, err
	}
	return animals, nil
}

func (s *AnimalStore) GetAllAnimals(c *gin.Context) ([]models.Animal, error) {
	animals := []models.Animal{}
	result := s.db.Scopes(s.addMediaPreload, s.buildPetQuery(c), pagination.Paginate(c)).Find(&animals)
	return animals, result.Error
}

func (s *AnimalStore) GetNotSeenAnimals(userID uuid.UUID, c *gin.Context) ([]models.Animal, error) {
	var animals []models.Animal
	twentyFourHoursAgo := time.Now().Add(-24 * time.Hour)

	err := s.db.
		Select("animals.*").
		Joins("LEFT JOIN seen_animals ON animals.id = seen_animals.animal_id AND seen_animals.user_id = ?", userID).
		Where("seen_animals.seen_at < ? OR seen_animals.seen_at IS NULL", twentyFourHoursAgo).
		Scopes(s.addMediaPreload, s.buildPetQuery(c), pagination.Paginate(c)).
		Find(&animals).Error
	if err != nil {
		return nil, err
	}

	return animals, nil
}

func (s *AnimalStore) MarkAsSeen(animalID uint, userID uuid.UUID, animalLiked bool) error {
	seenAnimal := models.SeenAnimal{AnimalID: animalID, UserID: userID, Liked: animalLiked, SeenAt: time.Now()}

	return s.db.Where("user_id = ? AND animal_id = ?", userID, animalID).Save(&seenAnimal).Error

}

func (s *AnimalStore) buildPetQuery(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		minAge := c.Query(constants.MinAgeParam)
		maxAge := c.Query(constants.MaxAgeParam)
		genders := c.Query(constants.GenderParam)
		location := c.Query(constants.LocationParam)
		vaccinated := c.Query(constants.VaccinatedParam)
		sterilized := c.Query(constants.SterilizedParam)

		if minAge != "" {
			db = db.Where("age >= ?", minAge)
			log.Info().Err(db.Error).Msg("age >=")
		}

		if maxAge != "" {
			db = db.Where("age <= ?", maxAge)
			log.Info().Err(db.Error).Msg("age <=")
		}

		if genders != "" {
			var g []string
			if err := json.Unmarshal([]byte(genders), &g); err == nil {
				db.Where("gender IN ?", g)
			}
		}

		if location != "" {
			db = db.Where("location = ?", location)
			log.Info().Err(db.Error).Any("location", location).Msg("location = ?")
		}

		if vaccinated != "" {
			v, _ := strconv.ParseBool(vaccinated)
			db = db.Where("vaccinated = ?", v)
			log.Info().Err(db.Error).Msg("vaccinated = ?")
		}

		if sterilized != "" {
			s, _ := strconv.ParseBool(sterilized)
			db = db.Where("sterilized = ?", s)
			log.Info().Err(db.Error).Msg("sterilized = ?")
		}

		return db
	}
}

func (s *AnimalStore) addMediaPreload(db *gorm.DB) *gorm.DB {
	return db.Preload("Photos").Preload("Image")
}
