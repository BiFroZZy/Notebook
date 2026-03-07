package auth

import (
	"context"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"

	l "notebook/internal/logger"
	mod "notebook/internal/models"
	db "notebook/internal/database"
)
var (
	ctx = context.Background()
	logger = l.NewLogger()
)
// @Summary Authorization posting data
// @Description Отправка данных пользователя в базу данных на странице авторизации
// @Router /public/auth/post [post]
// @Success 200
func Authorization(c echo.Context) error{
	conn, err := db.ConnectingSQL()
	if err != nil{
		logger.Err(err).Msg("Can't connect to database\n")
	}
	defer conn.Close(ctx)

	getAuthLogin := c.FormValue("auth_login")
	getAuthPassword := c.FormValue("auth_password")
	
	user := mod.User{}
	err = conn.QueryRow(ctx, os.Getenv("AUTH_QUERY"), getAuthLogin).Scan(&user.Login, &user.Password)
	if err != nil{
		logger.Err(err).Msg("Error in querying data in authorization\n")
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
		c.Redirect(http.StatusFound, "/users/"+userID+"/notes")
	} else {
		c.Render(http.StatusOK, "auth.html", map[string]interface{}{
			"Title": "Authorization",
			"Error": "Wrong login or password",
		})
	}
	return c.Redirect(http.StatusOK, "/public/reg")
}