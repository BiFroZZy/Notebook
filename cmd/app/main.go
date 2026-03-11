package main

import (
	"github.com/joho/godotenv"

	"notebook/cmd/handlers"
	l"notebook/internal/logger"
	_ "notebook/docs"
)

// @Title Notebook documentaion (Документация по программе)
// @Description Записки для пользователей
// @version 1.0
// @Host localhost
// @BasePath
func main() {
	logger := l.NewLogger()
	if err := godotenv.Load(); err != nil {
		logger.Fatal().Msg("Can't load env file\n")
	}
	logger.Info().Msg("App is started!")
	handlers.Handlers()
}
