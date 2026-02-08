package handlers

import (
	"net/http"
	"github.com/labstack/echo"
	// db "notebook/internal/database"
)
func MainPage(c echo.Context) error{
	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"Title": "Main",
		"Notes": "",
	})
}
func AuthPage(c echo.Context) error{
	return c.Render(http.StatusOK, "auth.html", map[string]interface{}{
		"Title": "Authorization",
	})
}
func RegPage(c echo.Context) error{
	return c.Render(http.StatusOK, "reg.html", map[string]interface{}{
		"Title": "Registration",
	})
}
func AboutPage(c echo.Context) error{
	return c.Render(http.StatusOK, "about.html", map[string]interface{}{
		"Title": "About",
	})
}