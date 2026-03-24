package jwt

import (
	"os"
	_"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"notebook/internal/config"
	l "notebook/internal/logger"
	mod "notebook/internal/models"
)

var (
	logger = l.NewLogger()
	secretKey = []byte(os.Getenv("SECRET_JWT"))
)

func GenerateJWT(userID uuid.UUID) (string, error) {
	cfg := config.Config{}
	mod.CallValidation(logger, cfg)
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"exp": cfg.Exp,
		"iat": cfg.Iat,
		"iss": cfg.Iss,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ParseToken(tokenString string) (*jwt.Token, error) {
    return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return secretKey, nil
    })
}