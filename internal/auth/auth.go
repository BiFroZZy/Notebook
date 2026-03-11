package auth

import (
	"context"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	db "notebook/internal/database"
	l "notebook/internal/logger"
	j "notebook/internal/jwt"
	mod "notebook/internal/models"
)
var (
	ctx = context.Background()
	logger = l.NewLogger()
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
		cookie, err := c.Cookie("token")
		if err != nil{
			logger.Err(err).Msg("Error occrued while gettting cookie\n")
			return c.Render(http.StatusOK, "auth.html", map[string]interface{}{
				"Title": "Authorization",
				"Error": "Your session is over, authorize again to continue",
			})
		}
		token, err := j.ParseToken(cookie.Value)
		if err != nil || !token.Valid {
			logger.Err(err).Msg("Error in parsing cookie value!\n")
			return c.Redirect(http.StatusUnauthorized, "/public/auth")
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("user_id", claims["user_id"])
		}
		return next(c)
    }
}

// @Summary Authorization posting data
// @Description Отправка данных пользователя в базу данных на странице авторизации
// @Router /public/auth/post [post]
// @Success 200
func Authorization(c echo.Context) error{
	conn, err := db.ConnectingSQL()
	if err != nil{
		logger.Err(err).Msg("Can't connect to database")
	}
	defer conn.Close(ctx)

	getAuthLogin := c.FormValue("auth_login")
	getAuthPassword := c.FormValue("auth_password")
	
	token, err := j.GenerateJWT(db.GetUserID())
	if err != nil{
		logger.Err(err).Msg("Error occrured while creating JWT")
	}
	// устанавливаю значение токена в куки, храня его в контексте
	c.SetCookie(&http.Cookie{
		Name: "token",
		Value: token,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure: false,
		Path: "/",
		MaxAge: 600,
	})

	user := mod.User{}

	err = conn.QueryRow(ctx, os.Getenv("AUTH_QUERY"), getAuthLogin).Scan(&user.Login, &user.Password)
	if err != nil{
		logger.Err(err).Msg("Error in querying data in authorization")
	}
	userID := db.GetUserID()

	if err == pgx.ErrNoRows{
		return c.Render(http.StatusOK, "auth.html", map[string]interface{}{
			"Title": "Authorization",
			"Error": "No such user!",
		})
	}
	ok := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(getAuthPassword))

	if ok == nil && user.Login == getAuthLogin{
		c.Redirect(http.StatusSeeOther, "/users/"+userID+"/notes")
	} 

	return c.Render(http.StatusOK, "auth.html", map[string]interface{}{
		"Title": "Authorization",
		"Error": "Wrong login or password",
	})
}