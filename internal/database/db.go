package database

import (
	"context"
	"log"
	"os"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo"
)

var ctx = context.Background()

func ConnectingSQL() (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, os.Getenv("PGX_URL"))
	if err != nil{
		log.Printf("Не могу подключиться к базе данных: %v\n", err)
	} 
	_, err = conn.Exec(ctx, `CREATE TABLE IF NOT EXISTS users (
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
	defer conn.Close(ctx)

	_, err = conn.Exec(ctx, 
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
	err = conn.QueryRow(ctx, "SELECT user_login, user_password, user_email FROM users WHERE user_login = $1 OR user_password = $2", getRegLogin, getRegPassword).Scan(&login, &password, &email)
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
	defer conn.Close(ctx)

	getAuthLogin := c.FormValue("auth_login")
	getAuthPassword := c.FormValue("auth_password")

	var login, password string
	err = conn.QueryRow(ctx, "SELECT user_login, user_password FROM users WHERE user_login = $1", getAuthLogin).Scan(&login, &password)
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

func GetUserID() string{
	conn, err := ConnectingSQL()
	if err != nil{
		log.Printf("Ошибка в получении ID: %v", err)
	}
	var userID string
	err = conn.QueryRow(ctx, "SELECT user_id FROM users").Scan(&userID)
	if err != nil{
		log.Printf("%v", err)
	}
	return userID
}


func WriteNotes(c echo.Context) error{
	notes := c.FormValue("write_notes")
	conn, err := ConnectingSQL()
	if err != nil{
		log.Printf("%v", err)
	}
	defer conn.Close(ctx)

	_, err = conn.Exec(ctx, "INSERT INTO users_notes(notes) VALUES ($1)", notes)
	if err != nil{
		log.Printf("Can't insert user's notes in DB: %v", err)
	}
	// data := map[string]interface{}{"Title": "Main", "Notes": notes}
	data := struct{Title string; Notes string}{Title: "Main", Notes: notes}
	// ages := struct{age int}{age: 42} Пример как записывать чтоб не заыбть как придурыч 
	return c.Render(http.StatusOK, "index.html", data)
}