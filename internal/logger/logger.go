package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/go-playground/validator/v10"

)
type Vals struct{
	*validator.Validate
}

func (vl Vals) ValidNew(logger zerolog.Logger, structure interface{}){
	if err := vl.Struct(structure); err != nil{
		logger.Error().Msg("Validation error occured!\n")
		for _, e := range err.(validator.ValidationErrors){
			logger.Err(err).Msg("Field: %s; Tag: %s"+ e.Field()+ e.Tag())
		}
	}
}

func CallValid(logger zerolog.Logger, structure interface{}){
	vals := Vals{}
	vals.ValidNew(logger, structure)
}

func NewLogger() zerolog.Logger{
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()
	return log.Logger
}