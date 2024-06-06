package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/Kachyr/findyourpet/findyourpet-backend/internal/store/users"
	"github.com/Kachyr/findyourpet/findyourpet-backend/pkg/auth"
	"github.com/Kachyr/findyourpet/findyourpet-backend/pkg/constants"
	"github.com/Kachyr/findyourpet/findyourpet-backend/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type UserHandlerInterface interface {
	SignUp(c *gin.Context)
	LogIn(c *gin.Context)
	GetUser(c *gin.Context)
	GetUserSettings(c *gin.Context)
	SetUserSettings(c *gin.Context)
	SingUp(c *gin.Context)
	getUserDataFromContext(c *gin.Context) (*models.User, error)
}

type UserHandler struct {
	authService auth.AuthServiceI
	store       users.UserStoreI
}

func NewUserHandler(auth auth.AuthServiceI, store users.UserStoreI) *UserHandler {
	return &UserHandler{authService: auth, store: store}
}

func (h *UserHandler) SignUp(c *gin.Context) {
	var body models.UserSingupJSON

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		log.Err(err).Msg("Failed to read body")
		return
	}

	hash, err := h.authService.GenerateHashFromPassword(body.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to hash password"})
		log.Err(err).Msg("Failed to hash password")
		return
	}

	user := models.User{
		Email:        body.Email,
		Password:     string(hash),
		UserSettings: constants.DefaultUserSettings,
	}
	err = h.store.Create(&user)

	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create user"})
		log.Err(err).Ctx(context.Background()).Msg("Failed to create user")
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *UserHandler) LogIn(c *gin.Context) {
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	user, err := h.store.GetByEmail(body.Email)
	if user.ID.ID() == 0 {
		log.Error().Err(err).Send()
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	tokenString, err := h.authService.Authenticate(*user, body.Password)
	if err != nil {
		log.Info().Err(err).Send()
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, int(constants.JWTExpireTime), "", "", false, true)
	c.JSON(http.StatusOK, gin.H{})
}

func (h *UserHandler) SetUserSettings(c *gin.Context) {
	user, err := h.getUserDataFromContext(c)

	if err != nil {
		c.Status(http.StatusBadRequest)
		log.Info().Err(err).Msg("Error getting user data")
		return
	}

	var body models.UserSettingsJSON

	if err := c.Bind(&body); err != nil {
		c.Status(http.StatusBadRequest)
		log.Info().Err(err).Msg("Error accessing request body")
		return
	}

	if err := h.store.SetUserSettings(user.ID, models.UserSettingsFromJSON(body)); err != nil {
		c.Status(http.StatusBadRequest)
		log.Info().Err(err).Msg("Error saving user settings")
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (h *UserHandler) GetUserSettings(c *gin.Context) {
	user, err := h.getUserDataFromContext(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		log.Info().Err(err).Msg("Error retrieving user settings")
		return
	}

	userSettings, err := h.store.GetUserSettings(user.ID)
	if err != nil {
		c.Status(http.StatusBadRequest)
		log.Info().Err(err).Msg("Error retrieving user settings")
		return
	}
	c.JSON(http.StatusOK, models.UserSettingsToJSON(userSettings))
}

func (h *UserHandler) GetUser(c *gin.Context) {
	user, err := h.getUserDataFromContext(c)

	if err != nil {
		c.Status(http.StatusUnauthorized)
		log.Info().Err(err).Msg("Error creating post")
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) getUserDataFromContext(c *gin.Context) (*models.User, error) {
	u, ok := c.Get("user")

	if !ok {
		return nil, errors.New("cannot get user data from cookie")
	}
	user, ok := u.(*models.User)
	if !ok {
		return nil, errors.New("cannot convert to user data from cookie")
	}
	return user, nil
}
