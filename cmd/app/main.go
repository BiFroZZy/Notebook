package main

import (
	"github.com/joho/godotenv"

	"notebook/cmd/server"
	_ "notebook/docs"
	"notebook/internal/config"
	l "notebook/internal/logger"
)

// @Title Notebook documentaion (Документация по программе)
// @Description Записки для пользователей
// @version 1.0
// @Host localhost
// @BasePath
func main() {
	logger := l.NewLogger()
	cfg, err := config.Load()
	if err != nil{
		logger.Fatal().Err(err).Msg("Can't load configs")
	}
	if err := godotenv.Load(); err != nil {
		logger.Fatal().Msg("Can't load env file")
	}
	srv := server.New(cfg)
	srv.Start()
	logger.Info().Msg("App is started!")
}
