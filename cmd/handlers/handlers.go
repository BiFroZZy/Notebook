package handlers

import (
	"html/template"
	"net/http"
	"io"
	"log"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	db "notebook/internal/database"
	h "notebook/web/handlers"
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

	e.Use(middleware.Logger())
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
	users := e.Group("/users")
	// users.GET("/user/:userID/main", MainPage)
	users.GET("/about", h.AboutPage)
	users.GET("/main", h.MainPage)
	users.POST("/main/post", db.WriteNotes)

	tmpl, err := template.ParseFiles(
		"web/templates/index.html", 
		"web/templates/auth.html",
		"web/templates/reg.html",
		"web/templates/header.html",
		"web/templates/footer.html",
		"web/templates/about.html",
	)
	if err != nil{
		log.Printf("Ошибка парсинга HTML-шаблона: %v\n", err)
	}

	e.Renderer = &Template{templates: tmpl}
	e.Static("/web/css/", "web/css/styles.css")

	e.Logger.Fatal(e.Start(":8070"))
}