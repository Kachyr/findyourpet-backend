package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type UserSettings struct {
	gorm.Model
	UserID     uuid.UUID
	Type       *string
	MinAge     int
	MaxAge     int
	Gender     pq.StringArray `gorm:"type:text[]"`
	Location   *string
	Vaccinated bool
	Sterilized bool
}

type UserSettingsJSON struct {
	UserID     uuid.UUID `json:"userId"`
	Type       *string   `json:"type"`
	MinAge     int       `json:"minAge"`
	MaxAge     int       `json:"maxAge"`
	Gender     []string  `json:"gender"`
	Location   *string   `json:"location"`
	Vaccinated bool      `json:"vaccinated"`
	Sterilized bool      `json:"sterilized"`
}

func UserSettingsFromJSON(settings UserSettingsJSON) UserSettings {

	result := UserSettings{
		UserID:     settings.UserID,
		Type:       settings.Type,
		MinAge:     settings.MinAge,
		MaxAge:     settings.MaxAge,
		Gender:     pq.StringArray(settings.Gender),
		Location:   settings.Location,
		Vaccinated: settings.Vaccinated,
		Sterilized: settings.Sterilized,
	}

	// Log JSON data for debugging
	// log.Info().Interface("input_settings", settings).Interface("resulting_settings", result).Msg("UserSettingsFromJSON")

	return result
}
func UserSettingsToJSON(settings UserSettings) UserSettingsJSON {

	result := UserSettingsJSON{
		UserID:     settings.UserID,
		Type:       settings.Type,
		MinAge:     settings.MinAge,
		MaxAge:     settings.MaxAge,
		Gender:     settings.Gender,
		Location:   settings.Location,
		Vaccinated: settings.Vaccinated,
		Sterilized: settings.Sterilized,
	}

	return result
}
