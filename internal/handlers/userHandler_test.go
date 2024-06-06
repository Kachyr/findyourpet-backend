package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Kachyr/findyourpet/findyourpet-backend/internal/handlers"
	"github.com/Kachyr/findyourpet/findyourpet-backend/mocks"
	"github.com/Kachyr/findyourpet/findyourpet-backend/pkg/constants"
	"github.com/Kachyr/findyourpet/findyourpet-backend/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_SignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mocks.NewMockAuthServiceI(ctrl)
	mockUserStore := mocks.NewMockUserStoreI(ctrl)
	userHandler := handlers.NewUserHandler(mockAuthService, mockUserStore)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	t.Run("Successful signup", func(t *testing.T) {
		hashedPassword := "hashed_password"
		expectedUser := models.User{
			Email:        "test@example.com",
			Password:     hashedPassword,
			UserSettings: constants.DefaultUserSettings,
		}

		mockAuthService.EXPECT().
			GenerateHashFromPassword("password").
			Return(hashedPassword, nil)

		mockUserStore.EXPECT().
			Create(&expectedUser).
			Return(nil)

		reqBody := models.UserSingupJSON{
			Email:    "test@example.com",
			Password: "password",
		}
		userJSON, _ := json.Marshal(reqBody)

		r, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(userJSON))
		r.Header.Set("Content-Type", "application/json")
		c.Request = r

		userHandler.SignUp(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

}

func TestUserHandler_LogIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mocks.NewMockAuthServiceI(ctrl)
	mockUserStore := mocks.NewMockUserStoreI(ctrl)
	userHandler := handlers.NewUserHandler(mockAuthService, mockUserStore)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	t.Run("Successful login", func(t *testing.T) {

		reqBody := models.UserSingupJSON{
			Email:    "test@example.com",
			Password: "password",
		}
		userJSON, _ := json.Marshal(reqBody)

		r, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(userJSON))
		r.Header.Set("Content-Type", "application/json")
		c.Request = r

		user := models.User{ID: uuid.New(), Email: "test@example.com", Password: "hashed_password"}
		mockUserStore.EXPECT().GetByEmail(reqBody.Email).Return(&user, nil)

		tokenStringMock := "someRandomToken"
		mockAuthService.EXPECT().Authenticate(user, reqBody.Password).Return(tokenStringMock, nil)

		userHandler.LogIn(c)
		cookie := w.Result().Cookies()[0]
		assert.Equal(t, "Authorization", cookie.Name)
		assert.Equal(t, tokenStringMock, cookie.Value)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestUserHandler_SetUserSettings(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mocks.NewMockAuthServiceI(ctrl)
	mockUserStore := mocks.NewMockUserStoreI(ctrl)
	userHandler := handlers.NewUserHandler(mockAuthService, mockUserStore)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	t.Run("Successful SetUserSettings", func(t *testing.T) {
		reqBody := models.UserSettingsJSON{
			Type:       stringPtr("dog"),
			MinAge:     1,
			MaxAge:     5,
			Gender:     []string{"MALE", "FEMALE"},
			Location:   stringPtr("New York"),
			Vaccinated: true,
			Sterilized: false,
		}
		userJSON, _ := json.Marshal(reqBody)

		r, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(userJSON))
		r.Header.Set("Content-Type", "application/json")
		c.Request = r
		userMock := &models.User{ID: uuid.New(), Email: "test123@email.com", Password: "hashPassword"}
		c.Set("user", userMock)

		expectedSettings := models.UserSettingsFromJSON(reqBody)

		mockUserStore.EXPECT().SetUserSettings(userMock.ID, expectedSettings).Return(nil)

		userHandler.SetUserSettings(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestUserHandler_GetUserSettings(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserStore := mocks.NewMockUserStoreI(ctrl)
	userHandler := handlers.NewUserHandler(nil, mockUserStore)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	t.Run("successful retrieval of user settings", func(t *testing.T) {
		userMock := &models.User{ID: uuid.New(), Email: "test123@email.com", Password: "hashed_password"}
		c.Set("user", userMock)

		expectedSettings := models.UserSettings{
			UserID:     userMock.ID,
			Type:       stringPtr("dog"),
			MinAge:     1,
			MaxAge:     5,
			Gender:     []string{"MALE", "FEMALE"},
			Location:   stringPtr("New York"),
			Vaccinated: true,
			Sterilized: false,
		}

		mockUserStore.EXPECT().GetUserSettings(userMock.ID).Return(expectedSettings, nil)

		r, _ := http.NewRequest("GET", "/settings", nil)
		c.Request = r

		userHandler.GetUserSettings(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var responseSettings models.UserSettingsJSON
		err := json.Unmarshal(w.Body.Bytes(), &responseSettings)
		assert.NoError(t, err)
		assert.Equal(t, models.UserSettingsToJSON(expectedSettings), responseSettings)
	})
}

func stringPtr(s string) *string {
	return &s
}
