package constants

import "github.com/Kachyr/findyourpet/findyourpet-backend/pkg/models"

const (
	MALE            = "MALE"
	FEMALE          = "FEMALE"
	DEFAULT_MIN_AGE = 0
	DEFAULT_MAX_AGE = 20
)

var DefaultUserSettings = models.UserSettings{
	MinAge:     DEFAULT_MIN_AGE,
	MaxAge:     DEFAULT_MAX_AGE,
	Gender:     []string{MALE, FEMALE},
	Location:   nil,
	Vaccinated: true,
	Sterilized: false,
}
