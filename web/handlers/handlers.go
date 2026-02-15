package handlers

import (
	"net/http"
	"github.com/labstack/echo/v4"
	// db "notebook/internal/database"
)
// func MainPage(c echo.Context) error{
// 	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
// 		"Title": "Main",
// 	})
// }

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

// @Summary About
// @Description Основная информация о проекте
// @Router /users/about [get]
// @Produce html
// @Success 200
func AboutPage(c echo.Context) error{
	return c.Render(http.StatusOK, "about.html", map[string]interface{}{
		"Title": "About",
	})
}