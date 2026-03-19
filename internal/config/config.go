package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/joho/godotenv"

	l "notebook/internal/logger"
	mod "notebook/internal/models"
)

var logger = l.NewLogger()

type Config struct{
	ServerPort 	string 	`envconfig:"SERVER_PORT" default:"8080"`
	PgxUrl 		string 	`envconfig:"PGX_URL" required:"true"`
	SecretJWT 	string	`envconfig:"SECRET_JWT" required:"true"`
}

func Load() (*Config, error){
	var cfg Config
	if err := godotenv.Load(); err != nil {
		logger.Fatal().Msg("Can't load env file")
	}
	if err := envconfig.Process("", &cfg); err != nil{
		logger.Error().Err(err)
		return nil, err
	}
	mod.CallValidation(logger, cfg)
	return &cfg, nil
}