package handlers

import (
	"errors"
	"net/http"

	"github.com/Kachyr/findyourpet/findyourpet-backend/internal/services"
	"github.com/Kachyr/findyourpet/findyourpet-backend/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type AnimalsHandler struct {
	animalService services.AnimalServiceI
}

func NewAnimalsHandler(animalService services.AnimalServiceI) *AnimalsHandler {
	return &AnimalsHandler{animalService: animalService}
}

func (h *AnimalsHandler) AddAnimal(c *gin.Context) {
	var body models.AnimalJSON
	if err := c.BindJSON(&body); err != nil {
		log.Info().Err(err).Send()
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := h.animalService.AddAnimal(&body); err != nil {
		log.Info().Err(err).Msg("Cant store animal record")
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{})

}

func (h *AnimalsHandler) GetAnimals(c *gin.Context) {
	user, err := getUserDataFromContext(c)
	if err != nil {
		log.Info().Err(err).Send()
		c.Status(http.StatusBadRequest)
		return
	}
	animals, err := h.animalService.GetAnimals(user.ID, c)
	if err != nil {
		log.Info().Err(err).Msg("Cant get animals")
		c.Status(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, models.PaginatedContent[models.AnimalJSON]{
		Data:       models.ToAnimalJSONArray(animals),
		Page:       c.GetInt("page"),
		PageSize:   c.GetInt("pageSize"),
		TotalPages: c.GetInt("totalPages"),
	})
}
func (h *AnimalsHandler) GetAllAnimals(c *gin.Context) {
	animals, err := h.animalService.GetAllAnimals(c)
	if err != nil {
		log.Info().Err(err).Msg("Cant get animals")
		c.Status(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, models.PaginatedContent[models.AnimalJSON]{
		Data:       models.ToAnimalJSONArray(animals),
		Page:       c.GetInt("page"),
		PageSize:   c.GetInt("pageSize"),
		TotalPages: c.GetInt("totalPages"),
	})
}

func (h *AnimalsHandler) GetLikedAnimals(c *gin.Context) {
	user, err := getUserDataFromContext(c)
	if err != nil {
		log.Info().Err(err).Send()
		c.Status(http.StatusBadRequest)
		return
	}

	animals, err := h.animalService.GetLikedAnimals(user.ID, c)
	if err != nil {
		log.Info().Err(err).Msg("Cant Liked get animals")
		c.Status(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, models.PaginatedContent[models.AnimalJSON]{
		Data:       models.ToAnimalJSONArray(animals),
		Page:       c.GetInt("page"),
		PageSize:   c.GetInt("pageSize"),
		TotalPages: c.GetInt("totalPages"),
	})
}

func (h *AnimalsHandler) GetAnimalByID(c *gin.Context) {
	animalId := c.Param("id")
	if animalId == "" {
		c.Status(http.StatusBadRequest)
		return
	}
	animal, err := h.animalService.GetAnimalById(animalId)
	if err != nil {
		log.Info().Err(err).Send()
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, models.ToAnimalJSON(animal))
}

func (h *AnimalsHandler) MarkAsSeen(c *gin.Context) {
	animalId := c.Param("id")
	if animalId == "" {
		errMsg := errors.New("missing animal id")
		log.Info().Err(errMsg).Send()
		c.String(http.StatusBadRequest, errMsg.Error())
		return
	}

	body := struct{ Like bool }{}
	if err := c.Bind(&body); err != nil {
		log.Info().Err(err).Send()
		c.Status(http.StatusBadRequest)
		return
	}

	user, err := getUserDataFromContext(c)
	if err != nil {
		log.Info().Err(err).Send()
		c.Status(http.StatusBadRequest)
		return
	}

	if err := h.animalService.MarkAsSeen(animalId, user.ID, body.Like); err != nil {
		log.Info().Err(err).Msg("Cant mark as seen")
		c.Status(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, gin.H{})

}

func getUserDataFromContext(c *gin.Context) (*models.User, error) {
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
