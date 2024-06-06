package db

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

func Connect(ctx context.Context, readConn, writeConn string, retries int, retryWaitDuration time.Duration) (*gorm.DB, error) {

	log.Ctx(ctx).Info().Bool("Write connection string is set", writeConn != "").Bool("Read connection string is set", readConn != "").Msg("Connecting to database")

	var err error
	var db *gorm.DB

	for i := retries; i > 0; i-- {
		log.Ctx(ctx).Info().Msgf("Try connect to db, retriesNumber=%d", i)

		db, err = gorm.Open(
			postgres.Open(writeConn),
			&gorm.Config{Logger: gormLogger.Default.LogMode(gormLogger.Silent), TranslateError: true},
		)
		if err == nil {
			err = db.Use(dbresolver.Register(dbresolver.Config{
				Replicas: []gorm.Dialector{postgres.Open(readConn)},
			}))
			if err == nil {
				break
			}
		}

		time.Sleep(retryWaitDuration)
	}

	err = db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error

	if err != nil {
		return nil, errors.Wrapf(err, "Failed to open a connection to Postgres database")
	}

	return db, nil
}
