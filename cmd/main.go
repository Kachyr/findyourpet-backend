package main

import (
	"context"
	"time"

	"github.com/Kachyr/findyourpet/findyourpet-backend/internal/initializers"
	"github.com/Kachyr/findyourpet/findyourpet-backend/internal/logger"
	"github.com/Kachyr/findyourpet/findyourpet-backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/Kachyr/findyourpet/findyourpet-backend/internal/store/animals"
	"github.com/Kachyr/findyourpet/findyourpet-backend/internal/store/users"
	"github.com/Kachyr/findyourpet/findyourpet-backend/pkg/awsS3"
	"github.com/Kachyr/findyourpet/findyourpet-backend/pkg/db"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

const appName = "findyourpet-backend"

var gormDB *gorm.DB
var configuration *config

func init() {
	// prepare config
	var err error
	configuration, err = loadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to prepare config")
	}
	log.Info().Str("ENV", viper.GetString(environment)).Err(err).Msg("config is ready")
	ctx := context.Background()

	initLogger(ctx, configuration.LogLevel)
	gormDB, err = connectDB(ctx, configuration.Database)

	if err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("unable to connect to database")
	}

	initializers.SyncDatabase(gormDB)
}

func main() {

	userStore := users.NewUserStore(gormDB)
	animalStore := animals.NewAnimalStore(gormDB)
	s3Service := awsS3.NewS3Service("findyourpet-kach")
	animalService := services.NewAnimalService(animalStore, s3Service)
	router := initializers.NewRouter(gormDB, userStore, animalStore, animalService)

	ginEngine := gin.Default()
	router.SetupAPIs(ginEngine)

	ginEngine.Run(configuration.GinPort)
}

func initLogger(ctx context.Context, logLevel zerolog.Level) {
	logger.Init(logLevel, appName)
	log.Ctx(ctx).Info().Msg("logger initialized")
}

func connectDB(ctx context.Context, dbConfig *dbConfig) (*gorm.DB, error) {
	// connect to database
	gormDB, err := db.Connect(ctx, dbConfig.ReadURL, dbConfig.WriteURL, 3, time.Second*2)

	return gormDB, err
}
