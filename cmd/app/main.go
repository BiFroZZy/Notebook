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
func RegPage(c echo.Context) error{
	return c.Render(http.StatusOK, "reg.html", map[string]interface{}{
		"Title": "Registration",
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
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS users(
		ID SERIAL PRIMARY KEY,
		user_id INTEGER,
		user_login VARCHAR(50),
		user_password VARCHAR(50),
		user_email VARCHAR(100)),
		created_at TIMESTAMP DEFAULT NOW()`)
	if err != nil{
		log.Printf("Can't create table: %v\n", err)
	}
	return conn, err
}

func WriteDataSQL(login, password, email string) {
	conn, err := ConnectingSQL()
	if err != nil {
		log.Printf("Database connection error: %v\n", err)
	}
	type User struct{
		ID int
		Login string
		Password string
		Email string
	}
	var user User
	err = conn.QueryRow(context.Background(), 
		`INSERT INTO users(user_id, user_login, user_password, user_email) 
		VALUES($1, $2, $3, $4)`).Scan(&user.ID, &user.Login, &user.Password, &user.Email)
	if err != nil{
		log.Printf("Can't insert data in table: %v\n", err)
	}
}		
func Authorization(c echo.Context) error{
	conn, err := ConnectingSQL()
	if err != nil{
		log.Printf("Database connection error: %v", err)
	}
	getAuthLogin := c.FormValue("auth_login")
	getAuthPassword := c.FormValue("auth_password")
	
	// getDBLogin := conn.QueryRow(context.Background(), "SELECT user_login FROM users")
	// getDBPassword := conn.QueryRow(context.Background(), "") 

	rows := conn.QueryRow(context.Background(), "SELECT user_login, user_password FROM users")
	type User struct{
		ID int
		Login string
		Password string
		Email string
	}
	var user User
	err = conn.QueryRow(context.Background(), 
		`SELECT users(user_id, user_login, user_password, user_email) 
		VALUES($1, $2, $3, $4)`).Scan(&user.ID, &user.Login, &user.Password, &user.Email)
	if err != nil{
		log.Printf("Can't insert data in table: %v\n", err)
	}
	
	return c.Redirect(http.StatusOK, "/reg")
}

func Handlers(){
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	
	e.GET("/", func(c echo.Context) error{
		return c.Render(http.StatusOK, "auth.html", nil)
	})
	public := e.Group("/public")
	public.GET("/auth", AuthPage)
	public.GET("/reg", RegPage)

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