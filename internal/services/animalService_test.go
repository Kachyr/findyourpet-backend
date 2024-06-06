package services_test

import (
	"testing"

	"github.com/Kachyr/findyourpet/findyourpet-backend/internal/services"
	"github.com/Kachyr/findyourpet/findyourpet-backend/mocks"
	"github.com/Kachyr/findyourpet/findyourpet-backend/pkg/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
)

func TestAnimalService_AddAnimal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAnimalStore := mocks.NewMockAnimalStoreI(ctrl)
	mockS3Service := mocks.NewMockS3ServiceI(ctrl)
	service := services.NewAnimalService(mockAnimalStore, mockS3Service)

	animalJSON := &models.AnimalJSON{
		Name:   "Test Animal",
		Image:  "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQABAAD/2wCEAAkGBxISE...",
		Photos: []string{"data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQABAAD/2wCEAAkGBxISE..."},
	}

	animal := models.FromAnimalJSON(animalJSON)
	uploadOutput := &manager.UploadOutput{
		Location: "https://s3.amazonaws.com/findyourpet-kach/Test_Animal_image.jpg",
		Key:      aws.String("Test_Animal_image.jpg"),
	}
	photoOutput := models.Photo{ImageURL: uploadOutput.Location, Key: *uploadOutput.Key}

	mockS3Service.EXPECT().UploadSinglePhoto(animalJSON.Image, animalJSON.Name+"_image").Return(uploadOutput, nil)
	mockS3Service.EXPECT().UploadPhotos(animalJSON.Photos, animalJSON.Name).Return([]models.Photo{photoOutput}, nil)

	expectedAnimal := animal
	expectedAnimal.Image = models.Image{
		URL: uploadOutput.Location,
		Key: *uploadOutput.Key,
	}
	expectedAnimal.Photos = []models.Photo{photoOutput}

	mockAnimalStore.EXPECT().AddAnimal(gomock.Any()).DoAndReturn(func(arg *models.Animal) error {
		assert.Equal(t, expectedAnimal.Name, arg.Name)
		assert.Equal(t, expectedAnimal.Image, arg.Image)
		assert.Equal(t, expectedAnimal.Photos, arg.Photos)
		return nil
	})

	err := service.AddAnimal(animalJSON)
	assert.NoError(t, err)
}

func TestAnimalService_GetAnimals(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAnimalStore := mocks.NewMockAnimalStoreI(ctrl)
	mockS3Service := mocks.NewMockS3ServiceI(ctrl)
	service := services.NewAnimalService(mockAnimalStore, mockS3Service)

	userID := uuid.New()
	ginContext := &gin.Context{}

	expectedAnimals := []models.Animal{
		{Name: "Animal 1"},
		{Name: "Animal 2"},
	}

	mockAnimalStore.EXPECT().GetNotSeenAnimals(userID, ginContext).Return(expectedAnimals, nil)

	animals, err := service.GetAnimals(userID, ginContext)
	assert.NoError(t, err)
	assert.Equal(t, expectedAnimals, animals)
}

func TestAnimalService_GetAllAnimals(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAnimalStore := mocks.NewMockAnimalStoreI(ctrl)
	mockS3Service := mocks.NewMockS3ServiceI(ctrl)
	service := services.NewAnimalService(mockAnimalStore, mockS3Service)

	ginContext := &gin.Context{}

	expectedAnimals := []models.Animal{
		{Name: "Animal 1"},
		{Name: "Animal 2"},
	}

	mockAnimalStore.EXPECT().GetAllAnimals(ginContext).Return(expectedAnimals, nil)

	animals, err := service.GetAllAnimals(ginContext)
	assert.NoError(t, err)
	assert.Equal(t, expectedAnimals, animals)
}

func TestAnimalService_GetAnimalById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAnimalStore := mocks.NewMockAnimalStoreI(ctrl)
	mockS3Service := mocks.NewMockS3ServiceI(ctrl)
	service := services.NewAnimalService(mockAnimalStore, mockS3Service)

	animalID := uuid.New().String()

	expectedAnimal := models.Animal{Name: "Animal 1"}

	mockAnimalStore.EXPECT().GetById(animalID).Return(expectedAnimal, nil)

	animal, err := service.GetAnimalById(animalID)
	assert.NoError(t, err)
	assert.Equal(t, expectedAnimal, animal)
}

func TestAnimalService_MarkAsSeen(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAnimalStore := mocks.NewMockAnimalStoreI(ctrl)
	mockS3Service := mocks.NewMockS3ServiceI(ctrl)
	service := services.NewAnimalService(mockAnimalStore, mockS3Service)

	animalID := "1"
	userID := uuid.New()
	like := true

	mockAnimalStore.EXPECT().MarkAsSeen(uint(1), userID, like).Return(nil)

	err := service.MarkAsSeen(animalID, userID, like)
	assert.NoError(t, err)
}

func TestAnimalService_GetLikedAnimals(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAnimalStore := mocks.NewMockAnimalStoreI(ctrl)
	mockS3Service := mocks.NewMockS3ServiceI(ctrl)
	service := services.NewAnimalService(mockAnimalStore, mockS3Service)

	userID := uuid.New()
	ginContext := &gin.Context{}

	expectedAnimals := []models.Animal{
		{Name: "Animal 1"},
		{Name: "Animal 2"},
	}

	mockAnimalStore.EXPECT().GetLikedAnimals(userID, ginContext).Return(expectedAnimals, nil)

	animals, err := service.GetLikedAnimals(userID, ginContext)
	assert.NoError(t, err)
	assert.Equal(t, expectedAnimals, animals)
}
