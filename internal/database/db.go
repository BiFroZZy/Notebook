package database

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"

	l "notebook/internal/logger"
	mod "notebook/internal/models"
)

var (
	ctx = context.Background()
	logger = l.NewLogger()
) 

func ConnectingSQL() (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, os.Getenv("PGX_URL"))
	if err != nil{
		logger.Err(err).Msg("Can't connect to database")
	} 
	
	_, err = conn.Exec(ctx, os.Getenv("CREATE_TABLE"))
	if err != nil{
		logger.Err(err).Msg("Can't create table")
	}
	return conn, err
}

func WriteDataSQL(id, login, password, email string){
	conn, err := ConnectingSQL()
	if err != nil {
		logger.Err(err).Msg("Can't connect to database while writing data")
	}
	defer conn.Close(ctx)
	_, err = conn.Exec(ctx, os.Getenv("WRITE_SQL_QUERY"), id, login, password, email)
	if err != nil{
		logger.Err(err).Msg("Can't insert data to database")
	}
}	

func GetUser(c echo.Context) (uuid.UUID, string, time.Time){
	conn, err := ConnectingSQL()
	if err != nil{
		logger.Err(err).Msg("Can't get an ID")
	}
	defer conn.Close(ctx)
	user := mod.User{}
	getUserLogin := c.FormValue("auth_login")

	err = conn.QueryRow(ctx, os.Getenv("GET_USER"), getUserLogin).Scan(&user.ID, &user.Email, &user.CreatedAt)
	if err != nil{
		if err == pgx.ErrNoRows{
			logger.Error().Msg("Nothing to get - rows are empty")
		}else{
			logger.Err(err).Msg("Error occured while querying the row")
		}
	}
	return user.ID, user.Email, user.CreatedAt
}

func GetNotes(c echo.Context) []mod.Note{
	conn, err := ConnectingSQL()
	if err != nil {
		logger.Err(err).Msg("Error occured while connecting to DB")
	}
	defer conn.Close(ctx)
	usersData := []mod.Note{} 
	note := mod.Note{}
	var t time.Time
	
	rows, err := conn.Query(ctx, os.Getenv("GET_NOTES"), c.Param("user_id"))
	if err != nil{
		logger.Err(err).Msg("Error occured while querying with rows")
	}
	defer rows.Close()

	for rows.Next(){
		err := rows.Scan(&note.NotesData, &t, &note.UUID)
		if err!= nil{
			if err == pgx.ErrNoRows{
				logger.Error().Msg("Notes are emtpy!")
			}
			logger.Err(err).Msg("Error occured while scaning data with rows")
		}
		note.CreatedAt = t.Format("2006-01-02 15:04")
		usersData = append(usersData, note)
	}
	return usersData
}

// @Summary User's notes here
// @Description Заметки пользователя
// @Router /users/:id/notes [get]
// @Produce html
// @Success 200
func ShowNotes(c echo.Context) error{
	notes := GetNotes(c)
	userID := c.Param("user_id") // НЕ УДАЛЯТЬ!!!!!!
	if userID == "" || userID == "00000000-0000-0000-0000-000000000000" {
        userID = c.FormValue("user_id")
    }
	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"Title": "Notes", 
		"Notes": notes,
		"UserID": userID,
	})
}

// @Summary Deleting user's notes here
// @Description Удаление заметок пользователя
// @Router /users/:id/notes/delete [post]
// @Redirection 303
func DeleteNotes(c echo.Context) error {
	conn, err := ConnectingSQL()
	if err != nil{
		logger.Err(err).Msg("Can't connect to database")
	}
	defer conn.Close(ctx)

	UUID := c.FormValue("note_id")
	userID := c.Param("user_id")
	
	_, err = conn.Exec(ctx, os.Getenv("DELETE_NOTES"), UUID)
	if err != nil{
		logger.Err(err).Msg("Can't delete the note")
	}
	return c.Redirect(http.StatusSeeOther, "/users/"+userID+"/notes")
}
// Запись заметок в БДшку
func WriteNotes(c echo.Context) error {
	notes := c.FormValue("write_notes")
	notesUUID := uuid.New()
	userID := c.Param("user_id")
	if notes != ""{
		conn, err := ConnectingSQL()
		if err != nil{
			logger.Err(err).Msg("Can't connect to database")
		}
		defer conn.Close(ctx)

		_, err = conn.Exec(ctx, os.Getenv("WRITE_NOTES"), notesUUID, notes, userID)
		if err != nil{
			logger.Err(err).Msg("Can't insert user's notes in DB")
		}
		return ShowNotes(c)
	}
	return nil
}