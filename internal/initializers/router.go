package initializers

import (
	"github.com/Kachyr/findyourpet/findyourpet-backend/internal/handlers"
	"github.com/Kachyr/findyourpet/findyourpet-backend/internal/services"
	"github.com/Kachyr/findyourpet/findyourpet-backend/internal/store/animals"
	"github.com/Kachyr/findyourpet/findyourpet-backend/internal/store/users"
	"github.com/Kachyr/findyourpet/findyourpet-backend/pkg/auth"
	"github.com/Kachyr/findyourpet/findyourpet-backend/pkg/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Router struct {
	db            *gorm.DB
	authService   auth.AuthServiceI
	userStore     *users.UserStore
	animalsStore  *animals.AnimalStore
	animalService *services.AnimalService
}

func NewRouter(db *gorm.DB, userStore *users.UserStore, animalsStore *animals.AnimalStore, animalService *services.AnimalService) *Router {
	authService := auth.NewAuthService()
	return &Router{
		db:            db,
		authService:   authService,
		userStore:     userStore,
		animalsStore:  animalsStore,
		animalService: animalService,
	}
}

func (r *Router) SetupAPIs(e *gin.Engine) {
	e.MaxMultipartMemory = 7 << 20 // 7 MiB
	r.setupUsers(e)
	r.setupAnimals(e)
}

func (r *Router) setupUsers(e *gin.Engine) {
	userController := handlers.NewUserHandler(r.authService, r.userStore)
	e.POST("/singup", userController.SignUp)
	e.POST("/login", userController.LogIn)
	e.GET("/user", middleware.RequireAuth(r.userStore), userController.GetUser)
	e.POST("/settings", middleware.RequireAuth(r.userStore), userController.SetUserSettings)
	e.GET("/settings", middleware.RequireAuth(r.userStore), userController.GetUserSettings)
}

func (r *Router) setupAnimals(e *gin.Engine) {
	animalsHandler := handlers.NewAnimalsHandler(r.animalService)
	e.POST("/animal", middleware.RequireAuth(r.userStore), animalsHandler.AddAnimal)
	e.GET("/animal/:id", middleware.RequireAuth(r.userStore), animalsHandler.GetAnimalByID)
	e.PUT("/markasseen/:id", middleware.RequireAuth(r.userStore), animalsHandler.MarkAsSeen)
	e.GET("/animal", middleware.RequireAuth(r.userStore), animalsHandler.GetAnimals)
	e.GET("/animal/all", middleware.RequireAuth(r.userStore), animalsHandler.GetAllAnimals)
	e.GET("/user/likes", middleware.RequireAuth(r.userStore), animalsHandler.GetLikedAnimals)
}
