package handlers

import (
	_"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/http-swagger"

	db "notebook/internal/database"
	h "notebook/web/handlers"
	l "notebook/internal/logger"
	_"notebook/docs"
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
	public.POST("/auth/post", db.Authorization)
	public.GET("/auth", h.AuthPage)
	public.GET("/reg", h.RegPage)
	public.POST("/reg/post", db.Registration)
	
// TODO: сделать users/user{uuid}, где uuid получается из базы данных
	
	// route := fmt.Sprintf("/notes/:%v", db.GetNoteID)

	users := e.Group("/users/:user_id")
	users.DELETE("/notes/:note_id/delete", db.DeleteNotes)
	users.GET("/about", h.AboutPage)
	users.GET("/notes", db.ShowNotes)
	users.POST("/notes/post", db.WriteNotes)

	tmpl, err := template.ParseGlob(
		"web/templates/*.html", 
	)
	if err != nil{
		logger.Err(err).Msg("Error in parsing HTML templates in handler!")
	}

	e.Renderer = &Template{templates: tmpl}
	e.Static("/web/css/", "web/css/styles.css")
	e.Logger.Fatal(e.Start(":9080"))
}