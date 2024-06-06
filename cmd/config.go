package main

import (
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

// .environment file variables
const (
	logLevelEnv = "LOG_LEVEL"
	// local db
	dbReadURL  = "DB_READ_URL"
	dbWriteURL = "DB_WRITE_URL"
	//
	// remote db
	dbReadURLrender  = "DB_RENDER_READ_EX_URL"
	dbWriteURLrender = "DB_RENDER_WRITE_EX_URL"
	//
	ginPortEnv  = "GIN_PORT"
	environment = "ENV"
)

const (
	dev  = "DEV"
	prod = "PROD"
)

type config struct {
	LogLevel zerolog.Level
	Database *dbConfig
	GinPort  string
}

type dbConfig struct {
	ReadURL  string
	WriteURL string
}

func loadConfig() (*config, error) {
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		// .env file does not exist
		viper.AutomaticEnv()
	} else {
		viper.SetConfigFile(".env")
		err := viper.ReadInConfig()
		if err != nil {
			return nil, errors.Wrap(err, "error reading config")
		}
	}

	logLevel, err := zerolog.ParseLevel(strings.ToLower(viper.GetString(logLevelEnv)))
	if err != nil {
		return nil, errors.Wrap(err, "error parsing log level")
	}

	if isDevEnv() {
		return &config{
			LogLevel: logLevel,
			Database: &dbConfig{
				ReadURL:  viper.GetString(dbReadURL),
				WriteURL: viper.GetString(dbWriteURL),
			},
			GinPort: ":" + viper.GetString(ginPortEnv),
		}, nil
	}
	if isProdEnv() {
		return &config{
			LogLevel: logLevel,
			Database: &dbConfig{
				ReadURL:  viper.GetString(dbReadURLrender),
				WriteURL: viper.GetString(dbWriteURLrender),
			},
			GinPort: ":" + viper.GetString(ginPortEnv),
		}, nil
	}
	return nil, errors.Wrap(err, "error reading config")
}

func isDevEnv() bool {
	return viper.GetString(environment) == dev
}

func isProdEnv() bool {
	return viper.GetString(environment) == prod
}
