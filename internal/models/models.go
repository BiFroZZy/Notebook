package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type JWT struct{
	UserID string		`validate:"required,uuid"`
	Exp int64			`validate:"required"`
	Iat int64			`validate:"required"`
	Iss string 			`validate:"required"`
}

type Note struct{
	ID string 			`validate:"required,uuid"`
	UUID uuid.UUID 		`validate:"required,uuid"`
	NotesData string 	`validate:"min=1,max=1000"`
	CreatedAt string 	`validate:"required"`
	UserID uuid.UUID 	`validate:"uuid"`
}

type User struct{
	ID uuid.UUID		`validate:"required,uuid"`
	Name string 		`validate:"required,min=2,max=40"`
	Login string 		`validate:"required,min=4,max=12"`
	Password string 	`validate:"required,min=4,max=12"`
	Email string 		`validate:"required,email"`
	CreatedAt time.Time `validate:"required"`
}

type ValidationStruct struct{
	*validator.Validate
}

func (vl ValidationStruct) NewValidation(logger zerolog.Logger, structure interface{}){
	if err := vl.Struct(structure); err != nil{
		logger.Err(err).Msg("Validation error occured!\n")
		for _, e := range err.(validator.ValidationErrors){
			logger.Err(err).Msg("Field: %s; Tag: %s"+ e.Field()+ e.Tag())
		}
	}
}

func CallValidation(logger zerolog.Logger, structure interface{}){
	valid := ValidationStruct{
		Validate: validator.New(),
	}
	valid.NewValidation(logger, structure)
}