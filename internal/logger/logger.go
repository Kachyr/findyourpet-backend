package logger

import (
	stdLog "log"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

const (
	appNameFieldName   = "appName"
	levelFieldName     = "level"
	messageFieldName   = "msg"
	timestampFieldName = "time"
	timeFieldFormat    = time.RFC3339
)

var (
	outputStream = os.Stdout
)

func Init(logLevel zerolog.Level, appName string) {
	zerolog.TimeFieldFormat = timeFieldFormat
	zerolog.LevelFieldName = levelFieldName
	zerolog.MessageFieldName = messageFieldName
	zerolog.TimestampFieldName = timestampFieldName

	zerolog.SetGlobalLevel(logLevel)

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	log.Logger = zerolog.New(outputStream).With().
		Timestamp().
		Str(appNameFieldName, appName).
		Logger()

	// Take over the default logger in case of use by 3rd-party libs
	stdLog.SetOutput(log.Logger)
}
