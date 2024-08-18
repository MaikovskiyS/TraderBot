package logger

import (
	"os"

	"github.com/rs/zerolog"
)

const timeFormat = "2006-01-02T15:04:05.000Z07:00"

func New(level zerolog.Level) zerolog.Logger {
	zerolog.TimeFieldFormat = timeFormat
	zerolog.SetGlobalLevel(level)

	return zerolog.New(os.Stderr).With().Timestamp().Logger()
}
