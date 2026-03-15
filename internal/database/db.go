package database

import (
	"context"
	"net/http"
	"os"
	
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"

	l "notebook/internal/logger"
	mod "notebook/internal/models"
	_"notebook/internal/jwt"
)

var (
	ctx = context.Background()
	validate *validator.Validate
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

func GetUserID(c echo.Context) uuid.UUID{
	conn, err := ConnectingSQL()
	if err != nil{
		logger.Err(err).Msg("Can't get an ID")
	}
	defer conn.Close(ctx)
	user := mod.User{}

	if err != nil{
		logger.Err(err).Msg("Error in cookie show notes!")
	}
	getUserLogin := c.FormValue("auth_login")

	err = conn.QueryRow(ctx, os.Getenv("GET_USER_ID"), getUserLogin).Scan(&user.ID)

	if err != nil{
		if err == pgx.ErrNoRows{
			logger.Error().Msg("Rows are empty")
		}else{
			logger.Err(err).Msg("Error in querying the row in getting user's ID")
		}
	}
	return user.ID
}

func GetNotes(c echo.Context) []mod.Note{
	conn, err := ConnectingSQL()
	if err != nil {
		logger.Err(err).Msg("Error in getting notes")
	}
	defer conn.Close(ctx)
	usersData := []mod.Note{} 
	note := mod.Note{}

	rows, err := conn.Query(ctx, os.Getenv("GET_NOTES"), c.Param("user_id"))
	if err != nil{
		logger.Err(err).Msg("Error in querying with rows")
	}
	defer rows.Close()

	for rows.Next(){
		if err := rows.Scan(&note.NotesData, &note.CreatedAt); err!= nil{
			logger.Err(err).Msg("Error in scaning data with rows")
		}
		if err == pgx.ErrNoRows{
			logger.Info().Msg("Notes are emtpy!")
		}
		usersData = append(usersData, note)
	}
	return usersData
}
// Возвращает ID заметки
func GetNoteID(c echo.Context) (string, uuid.UUID){
	conn, err := ConnectingSQL()
	if err != nil{
		logger.Err(err).Msg("Error in getting ID")
	}
	defer conn.Close(ctx)
	var (
		noteID string
		noteUUID uuid.UUID
	)
	noteInfo := GetNotes(c)
	for _, n := range noteInfo{
		if err := conn.QueryRow(ctx, os.Getenv("GET_NOTE_ID_UUID"), n.NotesData).Scan(&noteID, &noteUUID); err != nil{
			logger.Err(err).Msg("Error in querying the row in getting note id/uuid")
		}
	}
	return noteID, noteUUID
}

// @Summary User's notes here
// @Description Заметки пользователя
// @Router /users/:id/notes [get]
func ShowNotes(c echo.Context) error{
	info := GetNotes(c)
	userID := c.Param("user_id") // НЕ УДАЛЯТЬ!!!!!!
	// if userID == ""{
	// 	userID := GetUserID(c)
	// 	return c.Redirect(http.StatusSeeOther, "/users/"+userID.String()+"/notes")
	// }
	logger.Info().Msg(userID)
	_, UUID := GetNoteID(c)
	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"Title": "Notes", 
		"Notes": info,
		"UserID": userID,
		"NoteID": UUID,
	})
}
// Удаление заметок из БДшки
func DeleteNotes(c echo.Context) error {
	conn, err := ConnectingSQL()
	if err != nil{
		logger.Err(err).Msg("Can't connect to database")
	}
	defer conn.Close(ctx)
	
	// userID := c.Param("user_id")
	// из за того что UserID в index.html delete передается в виде структуры Notes - она пустая, но я хз
	logger.Debug().Msg("1."+ c.Request().Method)
	logger.Debug().Msg("2."+c.Request().URL.String())
	logger.Debug().Msg("3."+c.Path())
	logger.Debug().Msg("4."+c.Param("user_id"))

	_, noteUUID := GetNoteID(c)
	_, err = conn.Exec(ctx, os.Getenv("DELETE_NOTES"), noteUUID)
	if err != nil{
		logger.Err(err).Msg("Can't delete the note")
	}
	return ShowNotes(c)
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