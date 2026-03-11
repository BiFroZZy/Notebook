package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func NewLogger() zerolog.Logger{
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()
	return log.Logger
}