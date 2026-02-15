package database

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

var ctx = context.Background()

type User struct{
	NotesData string
	CreatedAt time.Time
}
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
// @Summary Registration posting data 
// @Description Отправка данных пользователя в базу данных на странице регистрации
// @Router /public/reg/post [post]
// @Success 200
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
	return c.Redirect(http.StatusFound, "/users/notes")
}	
// @Summary Authorization posting data
// @Description Отправка данных пользователя в базу данных на странице авторизации
// @Router /public/auth/post [post]
// @Success 200
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
		c.Redirect(http.StatusFound, "/users/notes")
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

func GetNoteID() string{
	conn, err := ConnectingSQL()
	if err != nil{
		log.Printf("Ошибка в получении ID: %v", err)
	}
	var noteID string
	if err := conn.QueryRow(ctx, "SELECT user_id FROM users").Scan(&noteID);err != nil{
		log.Printf("%v", err)
	}
	return noteID
}

// @Summary User's notes here
// @Description Заметки пользователя
// @Router /users/notes [get]
func ShowNotes(c echo.Context) error{
	info := WriteNotes(c)
	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"Title": "Notes", 
		"Notes": info,
	})
}
// Создать в базе данных ID для заметок - для /users/main/:id - чтобы удалить заметку с этой же id
func DeleteNotes(c echo.Context) error {
	conn, err := ConnectingSQL()
	if err != nil{
		log.Printf("%v", err)
	}
	noteID := c.Param("id")
	_, err = conn.Exec(ctx, "DELETE FROM users_notes WHERE id = $1", noteID)
	if err != nil{
		log.Printf("Can't delete the note: %v", err)
	}
	return c.Render(http.StatusOK, "index.html", ShowNotes)
}

func WriteNotes(c echo.Context) []User {
	notes := c.FormValue("write_notes")
	if notes != ""{
		conn, err := ConnectingSQL()
		if err != nil{
			log.Printf("%v", err)
		}
		defer conn.Close(ctx)

		_, err = conn.Exec(ctx, "INSERT INTO users_notes(notes) VALUES ($1)", notes)
		if err != nil{
			log.Printf("Can't insert user's notes in DB: %v", err)
		}
		
		usersData := []User{}
		rows, err := conn.Query(ctx, "SELECT notes, created_at FROM users_notes")
		if err != nil{
			log.Printf("Error in querying with rows: %v", err)
		}
		defer rows.Close()

		for rows.Next(){
			u := User{}
			if err := rows.Scan(&u.NotesData, &u.CreatedAt); err!= nil{
				log.Printf("Error in scaning data with rows: %v", err)
			}
			usersData = append(usersData, u)
		}
		// if notes == ""{
		// 	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		// 		"Title": "Main",
		// 		"Notes": usersData,
		// 		"Error": "Empty field!",
		// 	})
		// }
		// return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		// 	"Title": "Main", 
		// 	"Notes": usersData,
		// })
		return usersData
	}
	return nil
}