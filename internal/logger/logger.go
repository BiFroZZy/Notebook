package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/ziflex/lecho/v3"
)

func NewLogger() zerolog.Logger{
	log.Logger = zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
		FormatLevel: func(i interface{}) string {return strings.ToUpper(fmt.Sprintf("[%s]", i))},
		FormatMessage: func(i interface{}) string {return fmt.Sprintf("| %s |", i)},
		FormatCaller: func(i interface{}) string {return filepath.Base(fmt.Sprintf("%s", i))},
		PartsExclude: []string{zerolog.TimestampFieldName}}).With().Timestamp().Caller().Logger()
	return log.Logger
}

func LoggerMiddleawre() echo.MiddlewareFunc{
	logger := lecho.New(os.Stdout, lecho.WithTimestamp())
	return lecho.Middleware(lecho.Config{
		Logger:              logger,
		RequestLatencyLevel: zerolog.WarnLevel,        
		RequestLatencyLimit: 500 * time.Millisecond,  
	})
}