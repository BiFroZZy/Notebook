package jwt

import(
	"time"
	"context"
	"os"

	"github.com/golang-jwt/jwt/v5"

	l"notebook/internal/logger"
	mod "notebook/internal/models"
)

var (
	ctx = context.Background()
	logger = l.NewLogger()
	secretKey = []byte(os.Getenv("SECRET_JWT"))
)

func GenerateJWT(userID string) (string, error) {
	tokenData := mod.JWT{
		UserID: userID,
		Exp: time.Now().Add(time.Second*600).Unix(),
		Iat: time.Now().Unix(),
		Iss: "Vladimir Putin",
	}
	mod.CallValidation(logger, tokenData)
	claims := jwt.MapClaims{
		"sub": tokenData.UserID,
		"exp": tokenData.Exp,
		"iat": tokenData.Iat,
		"iss": tokenData.Iss,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ParseToken(tokenString string) (*jwt.Token, error) {
    return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return secretKey, nil
    })
}