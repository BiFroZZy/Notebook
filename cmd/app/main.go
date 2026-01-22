package main

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Template struct{
	templates *template.Template
}

// Метод для рендера шаблонов, обращается к структур Template
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error{
	return t.templates.ExecuteTemplate(w, name, data) 
}

func MainPage(c echo.Context) error{
	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"Title": "Main",
	})
}

func AuthPage(c echo.Context) error{
	return c.Render(http.StatusOK, "auth.html", map[string]interface{}{
		"Title": "Authorization",
	})
}
func ConnectingSQL() (*pgx.Conn, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", 
		os.Getenv("PGX_USERNAME"),
		os.Getenv("PGX_PASSWORD"),
		os.Getenv("PGX_HOST"),
		os.Getenv("PGX_PORT"),
		os.Getenv("PGX_DB"),
	)
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil{
		log.Printf("Не могу подключиться к базе данных: %v", err)
	}
	return conn, err
}

func WriteDataSQL(login, password, email string) {
	ConnectingSQL()
	// getRegUsername := c.FormValue("login")
	// getRegPassword := c.FormValue("password")

}	
func Authorization(c echo.Context){
	conn, err := ConnectingSQL()
	if err != nil{
		log.Printf("Database connection error: %v", err)
	}
	getAuthLogin := c.FormValue("auth_login")
	getAuthPassword := c.FormValue("password")
	
	getDatabaseLogin = conn.QueryRow(context.Background(), "")
	getDatabasePassword = conn.QueryRow(context.Background(), "") 

	
}

func Handlers(){
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	
	e.GET("/", func(c echo.Context) error{
		return c.Render(http.StatusOK, "auth.html", nil)
	})
	e.GET("/auth", AuthPage)

	users := e.Group("/users")
	users.GET("/main", MainPage)
	
	tmpl, err := template.ParseFiles(
		"web/templates/index.html", 
		"web/templates/auth.html",
		"web/templates/header.html",
		"web/templates/footer.html",
)
	if err != nil{
		log.Printf("Ошибка парсинга HTML-шаблона: %v", err)
	}

	e.Renderer = &Template{templates: tmpl}
	e.Static("/web/css/", "web/css/styles.css")

	e.Logger.Fatal(e.Start(":8000"))
}

func main(){
	Handlers()
}