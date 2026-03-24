package handlers

import (
	"context"
	// "time"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"

	db "notebook/internal/database"
	l "notebook/internal/logger"
	mod "notebook/internal/models"
)
var(
	logger = l.NewLogger()
	ctx = context.Background()
)
// @Summary Authorization 
// @Description Страница авторизации
// @Router /public/auth [get]
// @Produce html
// @Success 200
func AuthPage(c echo.Context) error{
	return c.Render(http.StatusOK, "auth.html", map[string]interface{}{
		"Title": "Authorization",
	})
}
// @Summary Registration 
// @Description Страница регистрации
// @Router /public/reg [get]
// @Produce html
// @Success 200
func RegPage(c echo.Context) error{
	return c.Render(http.StatusOK, "reg.html", map[string]interface{}{
		"Title": "Registration",
	})
}

// @Summary Info
// @Description Основная информация о проекте
// @Router /users/:id/info [get]
// @Produce html
// @Success 200
func UserInfoPage(c echo.Context) error{
	userID := c.Param("user_id")
	user := mod.User{}
	conn, err := db.ConnectingSQL()
	if err != nil{
		logger.Err(err).Msg("Error")
	}
	defer conn.Close(ctx)
	
	err = conn.QueryRow(ctx, os.Getenv("GET_INFO"), userID).Scan(&user.CreatedAt, &user.Email)
	if err != nil{
		logger.Err(err).Msg("error")
	}

	return c.Render(http.StatusOK, "info.html", map[string]interface{}{
		"Title": "Info",
		"TimeInfo": user.CreatedAt.Format("2006-01-02 15:04:05"),
		"EmailInfo": user.Email,
		"UserID": userID,
	})
}	