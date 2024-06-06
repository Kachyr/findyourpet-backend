package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Kachyr/findyourpet/findyourpet-backend/internal/handlers"
	"github.com/Kachyr/findyourpet/findyourpet-backend/mocks"
	"github.com/Kachyr/findyourpet/findyourpet-backend/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAnimalsHandler_AddAnimal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	animalServiceMock := mocks.NewMockAnimalServiceI(ctrl)
	animalsHandler := handlers.NewAnimalsHandler(animalServiceMock)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	t.Run("Successful AddAnimal", func(t *testing.T) {
		reqBody := models.AnimalJSON{Name: "TEST", ID: 1, Age: 1, Type: "cat", Description: "qwerty", Gender: "MALE"}

		userJSON, _ := json.Marshal(reqBody)

		r, _ := http.NewRequest("POST", "/animal", bytes.NewBuffer(userJSON))
		r.Header.Set("Content-Type", "application/json")
		c.Request = r

		animalServiceMock.EXPECT().AddAnimal(&reqBody).Return(nil)

		animalsHandler.AddAnimal(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestAnimalsHandler_GetAnimals(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	animalServiceMock := mocks.NewMockAnimalServiceI(ctrl)
	animalsHandler := handlers.NewAnimalsHandler(animalServiceMock)

	t.Run("Successful GetAnimals", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		userMock := &models.User{ID: uuid.New(), Email: "test123@email.com", Password: "hashed_password"}
		c.Set("user", userMock)

		expectedAnimals := []models.Animal{
			{Name: "Animal1"},
			{Name: "Animal2"},
		}

		animalServiceMock.EXPECT().GetAnimals(userMock.ID, c).Return(expectedAnimals, nil)

		r, _ := http.NewRequest("GET", "/animals", nil)
		c.Request = r
		c.Set("page", 1)
		c.Set("pageSize", 10)
		c.Set("totalPages", 1)

		animalsHandler.GetAnimals(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response models.PaginatedContent[models.AnimalJSON]
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, models.ToAnimalJSONArray(expectedAnimals), response.Data)
		assert.Equal(t, 1, response.Page)
		assert.Equal(t, 10, response.PageSize)
		assert.Equal(t, 1, response.TotalPages)
	})

}

func TestAnimalsHandler_MarkAsSeen(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	animalServiceMock := mocks.NewMockAnimalServiceI(ctrl)
	animalsHandler := handlers.NewAnimalsHandler(animalServiceMock)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	t.Run("Successful MarkAsSeen", func(t *testing.T) {
		userMock := &models.User{ID: uuid.New(), Email: "test123@email.com", Password: "hashed_password"}
		c.Set("user", userMock)

		reqBody := struct{ Like bool }{Like: true}
		body, _ := json.Marshal(reqBody)

		r, _ := http.NewRequest("POST", "/animals/1/seen", bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")
		c.Request = r
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		animalServiceMock.EXPECT().MarkAsSeen("1", userMock.ID, true).Return(nil)

		animalsHandler.MarkAsSeen(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
