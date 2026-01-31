package main

import (
	"context"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
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
	conn, err := pgx.Connect(context.Background(), os.Getenv("PGX_URL"))
	if err != nil{
		log.Printf("Не могу подключиться к базе данных: %v\n", err)
	} 
	_, err = conn.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS users (
		ID SERIAL PRIMARY KEY,
		user_id VARCHAR(100),
		user_login VARCHAR(50),
		user_password VARCHAR(50),
		user_email VARCHAR(100),
		created_at TIMESTAMP DEFAULT NOW())`)
	if err != nil{
		log.Printf("Can't create table: %v\n", err)
	}
	return conn, err
}

func WriteDataSQL(id, login, password, email string){
	conn, err := ConnectingSQL()
	if err != nil {
		log.Printf("Database connection error: %v\n", err)
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(), 
		`INSERT INTO users (user_id, user_login, user_password, user_email) 
		VALUES ($1, $2, $3, $4)`, id, login, password, email)
	if err != nil{
		log.Printf("Can't insert data in table: %v\n", err)
	}
}		
func Registration(c echo.Context) error {
	getRegLogin := c.FormValue("reg_login")
	getRegPassword := c.FormValue("reg_password")
	getRegEmail := c.FormValue("reg_email")
	newUUID := uuid.New()
	
	conn, err := ConnectingSQL()
	if err != nil {
		return err
	}
	var login, password, email string
	err = conn.QueryRow(context.Background(), "SELECT user_login, user_password, user_email FROM users WHERE user_login = $1 OR user_password = $2", getRegLogin, getRegPassword).Scan(&login, &password, &email)
	if err != nil{
		return err
	}
	if getRegLogin == login && getRegPassword == password{
		return c.Render(http.StatusOK, "reg.html", map[string]interface{}{
			"Title": "Registration",
			"Error": "Login or password already exist",
		})
	}
	if getRegEmail == email{
		return c.Render(http.StatusOK, "reg.html", map[string]interface{}{
			"Title": "Registration",
			"Error": "You already have an account!",
		})
	}	
	WriteDataSQL(newUUID.String(), getRegLogin, getRegPassword, getRegEmail)
	return c.Redirect(http.StatusFound, "/users/main")
}	

func Authorization(c echo.Context) error{
	conn, err := ConnectingSQL()
	if err != nil{
		log.Printf("Database connection error: %v\n", err)
	}
	defer conn.Close(context.Background())

	getAuthLogin := c.FormValue("auth_login")
	getAuthPassword := c.FormValue("auth_password")

	var login, password string
	err = conn.QueryRow(context.Background(), "SELECT user_login, user_password FROM users WHERE user_login = $1", getAuthLogin).Scan(&login, &password)
	if err != nil{
		log.Printf("%v", err)
	}
	if err == pgx.ErrNoRows{
		return c.Render(http.StatusOK, "auth.html", map[string]interface{}{
			"Title": "Authorization",
			"Error": "No such user!",
		})
	}
	if password == getAuthPassword && login == getAuthLogin{
		c.Redirect(http.StatusFound, "/users/main")
	} else {
		c.Render(http.StatusOK, "auth.html", map[string]interface{}{
			"Title": "Authorization",
			"Error": "Wrong login or password",
		})
	}
	return c.Redirect(http.StatusOK, "/public/reg")
}

func Handlers(){
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	public := e.Group("/public")
	public.GET("/", func(c echo.Context) error{
		return c.Render(http.StatusOK, "auth.html", nil)
	})
	public.GET("/auth", AuthPage)
	public.POST("/auth/post", Authorization)
	public.GET("/reg", RegPage)
	public.POST("/reg/post", Registration)
// TODO: сделать users/user{uuid}, где uuid получается из базы данных
	users := e.Group("/users")
	users.GET("/main", MainPage)

	tmpl, err := template.ParseFiles(
		"web/templates/index.html", 
		"web/templates/auth.html",
		"web/templates/reg.html",
		"web/templates/header.html",
		"web/templates/footer.html",
	)
	if err != nil{
		log.Printf("Ошибка парсинга HTML-шаблона: %v\n", err)
	}

	e.Renderer = &Template{templates: tmpl}
	e.Static("/web/css/", "web/css/styles.css")

	e.Logger.Fatal(e.Start(":8080"))
}

func main(){
	if err := godotenv.Load(); err != nil{
		log.Fatalf("Can't load env file: %v\n", err)
	}
	Handlers()
}