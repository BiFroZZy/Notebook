package handlers

import (
	_ "fmt"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/http-swagger"

	_ "notebook/docs"
	"notebook/internal/auth"
	db "notebook/internal/database"
	l "notebook/internal/logger"
	h "notebook/web/handlers"
)

var logger = l.NewLogger()

type Template struct{
	templates *template.Template
}

// Метод для рендера шаблонов, обращается к структур Template
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error{
	return t.templates.ExecuteTemplate(w, name, data) 
}

func Handlers(){
	e := echo.New()
	
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Secure())

	e.GET("/swagger/*", echo.WrapHandler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
        httpSwagger.DocExpansion("none"),
	)))
	e.GET("/", func(c echo.Context) error{
		return c.Render(http.StatusOK, "auth.html", nil)
	})

	public := e.Group("/public")
	{	
		public.POST("/auth/post", auth.Authorization)
		public.GET("/auth", h.AuthPage)
		public.GET("/reg", h.RegPage)
		public.POST("/reg/post", auth.Registration)
	}
	users := e.Group("/users/:user_id")
	users.Use(auth.AuthMiddleware)
	{
		users.GET("/about", h.AboutPage)
		users.GET("/notes", db.ShowNotes)
		users.POST("/notes/delete", db.DeleteNotes)
		users.POST("/notes/post", db.WriteNotes)
	}

	tmpl, err := template.ParseGlob(
		"web/templates/*.html", 
	)
	if err != nil{
		logger.Err(err).Msg("Error in parsing HTML templates in handler!")
	}

	e.Renderer = &Template{templates: tmpl}
	e.Logger.Fatal(e.Start(os.Getenv("SERVER_PORT")))
}