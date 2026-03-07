package main

import (
	//"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"

	"notebook/cmd/handlers"
	_ "notebook/docs"
	l"notebook/internal/logger"
)

// var (
// 	validate *validator.Validate
// )

// @Title Notebook documentaion (Документация по программе)
// @Description Записки для пользователей
// @version 1.0
// @Host localhost
// @BasePath
func main() {
	logger := l.NewLogger()
	//validate := validator.New()
	if err := godotenv.Load(); err != nil {
		logger.Fatal().Msg("Can't load env file\n")
	}
	logger.Info().Msg("App is started!")
	handlers.Handlers()
}
