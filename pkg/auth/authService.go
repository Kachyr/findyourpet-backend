package auth

import (
	"errors"
	"os"
	"time"

	"github.com/Kachyr/findyourpet/findyourpet-backend/pkg/constants"
	"github.com/Kachyr/findyourpet/findyourpet-backend/pkg/models"
	"github.com/golang-jwt/jwt/v5"

	"golang.org/x/crypto/bcrypt"
)

type AuthServiceI interface {
	Authenticate(user models.User, userPassword string) (string, error)
	GenerateHashFromPassword(password string) (string, error)
}

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (a *AuthService) Authenticate(user models.User, userPassword string) (string, error) {
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userPassword)) != nil {
		return "", errors.New("wrong password")
	}

	// Initialize jwt.MapClaims as an empty map
	claims := jwt.MapClaims{}
	claims["sub"] = user.ID
	claims["exp"] = time.Now().Add(constants.CookieLifetime).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	return tokenString, err
}

func (a *AuthService) GenerateHashFromPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hash), err
}
