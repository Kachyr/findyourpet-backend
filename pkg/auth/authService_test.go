package auth_test

import (
	"os"
	"testing"
	"time"

	"github.com/Kachyr/findyourpet/findyourpet-backend/pkg/auth"
	"github.com/Kachyr/findyourpet/findyourpet-backend/pkg/constants"
	"github.com/Kachyr/findyourpet/findyourpet-backend/pkg/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Authenticate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a real instance of AuthService
	authService := auth.NewAuthService()

	t.Run("successful authentication", func(t *testing.T) {
		user := models.User{
			ID:       uuid.New(),
			Password: "hashedPassword",
		}
		userPassword := "password"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(userPassword), 10)
		user.Password = string(hashedPassword)

		// Call the method being tested
		resultTokenString, err := authService.Authenticate(user, userPassword)

		// Assert the expected results
		assert.NoError(t, err)
		assert.NotEmpty(t, resultTokenString)
		println(resultTokenString)

		// Validate the token
		token, err := jwt.Parse(resultTokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET")), nil
		})
		assert.NoError(t, err)
		assert.True(t, token.Valid)

		claims, ok := token.Claims.(jwt.MapClaims)
		assert.True(t, ok)
		assert.Equal(t, user.ID.String(), claims["sub"])

		exp, ok := claims["exp"].(float64)
		assert.True(t, ok)
		assert.WithinDuration(t, time.Now().Add(constants.CookieLifetime), time.Unix(int64(exp), 0), time.Second)
	})

	t.Run("wrong password", func(t *testing.T) {
		user := models.User{
			ID:       uuid.New(),
			Password: "hashedPassword",
		}
		userPassword := "wrongPassword"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), 10)
		user.Password = string(hashedPassword)

		// Call the method being tested
		resultTokenString, err := authService.Authenticate(user, userPassword)

		// Assert the expected results
		assert.Empty(t, resultTokenString)
		assert.EqualError(t, err, "wrong password")
	})
}

func TestAuthService_GenerateHashFromPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authService := auth.NewAuthService()

	t.Run("successful hash generation", func(t *testing.T) {
		password := "password"

		resultHash, err := authService.GenerateHashFromPassword(password)

		assert.NoError(t, err)
		assert.NotEmpty(t, resultHash)

		err = bcrypt.CompareHashAndPassword([]byte(resultHash), []byte(password))
		assert.NoError(t, err)
	})
}
