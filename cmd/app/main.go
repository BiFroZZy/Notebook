package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"

	"notebook/cmd/handlers"
	_ "notebook/docs"
	db "notebook/internal/database"
	"notebook/internal/logger"
)

// @Title Notebook documentaion (Документация по программе)
// @Description Записки для пользователей
// @version 1.0
// @Host localhost
// @BasePath
var (
	validate *validator.Validate
	note = db.Note{}
)
func main() {
	logger := logger.NewLogger()
	validate = validator.New()
	if err := godotenv.Load(); err != nil {
		logger.Fatal().Msg("Can't load env file\n")
	}
	if err := validate.Struct(note); err != nil{
		logger.Error().Msg("Validation error occured!\n")
	}
	logger.Info().Msg("App is started!")
	handlers.Handlers()
}
