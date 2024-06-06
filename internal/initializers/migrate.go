package initializers

import (
	"github.com/Kachyr/findyourpet/findyourpet-backend/pkg/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func SyncDatabase(db *gorm.DB) {
	log.Info().Msg("Syncing database")

	if err := db.AutoMigrate(
		&models.User{},
		&models.UserSettings{},
		&models.Image{},
		&models.Photo{},
		&models.Animal{},
		&models.SeenAnimal{},
		// &models.LikedAnimal{},
	); err != nil {
		log.Fatal().Err(err).Msg("Error to migrate database")
	}
	if err := db.SetupJoinTable(&models.User{}, "SeenAnimals", &models.SeenAnimal{}); err != nil {
		log.Fatal().Err(err).Msg("Error to setup join table SeenAnimals")
	}
	// if err := db.SetupJoinTable(&models.User{}, "LikedAnimals", &models.LikedAnimal{}); err != nil {
	// 	log.Fatal().Err(err).Msg("Error to setup join table LikedAnimals")
	// }
}
