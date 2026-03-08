package models 

import (
	"time"
	
	"github.com/google/uuid"
)

type Note struct{
	ID string 			`validate:"required,uuid"`
	NoteUUID uuid.UUID 		`validate:"required,uuid"`
	NotesData string 	`validate:"min=1,max=1000"`
	CreatedAt time.Time `validate:"required"`
	UserID string 		`validate:"uuid"`
	NoteID string 		`validate:"uuid"`
}

type User struct{
	ID string			`validate:"required"`
	Name string 		`validate:"required,min=2,max=40"`
	Login string 		`validate:"required,min=4,max=12"`
	Password string 	`validate:"required,min=4,max=12"`
	Email string 		`validate:"required,email"`
}
