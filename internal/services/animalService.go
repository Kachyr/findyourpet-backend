package services

import (
	"strconv"

	"github.com/Kachyr/findyourpet/findyourpet-backend/internal/store/animals"
	"github.com/Kachyr/findyourpet/findyourpet-backend/pkg/awsS3"
	"github.com/Kachyr/findyourpet/findyourpet-backend/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AnimalServiceI interface {
	AddAnimal(animal *models.AnimalJSON) error
	GetAllAnimals(c *gin.Context) ([]models.Animal, error)
	GetAnimalById(id string) (models.Animal, error)
	GetAnimals(id uuid.UUID, c *gin.Context) ([]models.Animal, error)
	GetLikedAnimals(userID uuid.UUID, c *gin.Context) ([]models.Animal, error)
	MarkAsSeen(animalID string, userID uuid.UUID, like bool) error
}

type AnimalService struct {
	s3Service   awsS3.S3ServiceI
	animalStore animals.AnimalStoreI
}

func NewAnimalService(animalStore animals.AnimalStoreI, s3Service awsS3.S3ServiceI) *AnimalService {
	return &AnimalService{
		animalStore: animalStore,
		s3Service:   s3Service,
	}
}

func (s *AnimalService) GetAnimals(id uuid.UUID, c *gin.Context) ([]models.Animal, error) {
	return s.animalStore.GetNotSeenAnimals(id, c)
}

func (s *AnimalService) GetAllAnimals(c *gin.Context) ([]models.Animal, error) {
	return s.animalStore.GetAllAnimals(c)
}

func (s *AnimalService) GetAnimalById(id string) (models.Animal, error) {
	return s.animalStore.GetById(id)
}

func (s *AnimalService) AddAnimal(animal *models.AnimalJSON) error {
	a := models.FromAnimalJSON(animal)
	result, err := s.s3Service.UploadSinglePhoto(animal.Image, animal.Name+"_image")
	if err != nil {
		return err
	}
	a.Image = models.Image{URL: result.Location, Key: *result.Key}

	photosURLs, err := s.s3Service.UploadPhotos(animal.Photos, animal.Name)
	if err != nil {
		return err
	}
	a.Photos = photosURLs

	if err := s.animalStore.AddAnimal(a); err != nil {
		return err
	}
	return nil
}

func (s *AnimalService) MarkAsSeen(animalID string, userID uuid.UUID, like bool) error {
	aID, err := strconv.ParseUint(animalID, 10, 64)
	if err != nil {
		return err
	}

	return s.animalStore.MarkAsSeen(uint(aID), userID, like)
}

func (s *AnimalService) GetLikedAnimals(userID uuid.UUID, c *gin.Context) ([]models.Animal, error) {
	return s.animalStore.GetLikedAnimals(userID, c)
}
