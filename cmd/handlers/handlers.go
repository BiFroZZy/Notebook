package handlers

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/http-swagger"

	db "notebook/internal/database"
	h "notebook/web/handlers"
	_"notebook/docs"
)

type Template struct{
	templates *template.Template
}

// Метод для рендера шаблонов, обращается к структур Template
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error{
	return t.templates.ExecuteTemplate(w, name, data) 
}

func Handlers(){
	e := echo.New()

	e.GET("/swagger/*", echo.WrapHandler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
        httpSwagger.DocExpansion("none"),
	)))
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Secure())
	public := e.Group("/public")

	e.GET("/", func(c echo.Context) error{
		return c.Render(http.StatusOK, "auth.html", nil)
	})
	public.POST("/auth/post", db.Authorization)
	public.GET("/auth", h.AuthPage)
	public.GET("/reg", h.RegPage)
	public.POST("/reg/post", db.Registration)
	
// TODO: сделать users/user{uuid}, где uuid получается из базы данных
	noteID := db.GetNoteID()
	route := fmt.Sprintf("/notes/:%s", noteID)
	users := e.Group("/users")
	users.DELETE(route, db.DeleteNotes)
	users.GET("/about", h.AboutPage)
	users.GET("/notes", db.ShowNotes)
	users.POST("/notes/post", db.ShowNotes)
	tmpl, err := template.ParseGlob(
		"web/templates/*.html", 
	)
	if err != nil{
		log.Printf("Ошибка парсинга HTML-шаблона: %v\n", err)
	}

	e.Renderer = &Template{templates: tmpl}
	e.Static("/web/css/", "web/css/styles.css")
	
	e.Logger.Fatal(e.Start(":8080"))
}