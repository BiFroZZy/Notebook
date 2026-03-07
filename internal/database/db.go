package database

import (
	"context"
	"net/http"
	"os"
	
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	l "notebook/internal/logger"
	mod "notebook/internal/models"
)

var (
	ctx = context.Background()
	validate *validator.Validate
	logger = l.NewLogger()
) 

func HashingFunc(password string) (hashedPass []byte){
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil{
		logger.Err(err).Msg("Error while hashing the password!")
	}
	return hashedPass
}

func ConnectingSQL() (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, os.Getenv("PGX_URL"))
	if err != nil{
		logger.Err(err).Msg("Can't connect to database\n")
	} 
	_, err = conn.Exec(ctx, os.Getenv("CREATE_TABLE"))
	if err != nil{
		logger.Err(err).Msg("Can't create table\n")
	}
	return conn, err
}

func WriteDataSQL(id, login, password, email string){
	conn, err := ConnectingSQL()
	if err != nil {
		logger.Err(err).Msg("Can't connect to database while writing data\n")
	}
	defer conn.Close(ctx)

	_, err = conn.Exec(ctx, os.Getenv("WRITE_SQL_QUERY"), id, login, password, email)
	if err != nil{
		logger.Err(err).Msg("Can't insert data to database\n")
	}
}	

// пофиксить баг где id берется не правильно, по почте, почему то из базы даннных берет самого ластового чела после его регистрации
func GetUserID() string{
	conn, err := ConnectingSQL()
	if err != nil{
		logger.Err(err).Msg("Can't get an ID\n")
	}
	var userID, email string
	rows, err := conn.Query(ctx, os.Getenv("GET_USER_ID_EMAIL"))
	if err != nil{
		logger.Err(err).Msg("Error in querying userID\n")
	}
	defer rows.Close()

	for rows.Next(){
		err = rows.Scan(&email)
		if err != nil{
			logger.Err(err).Msg("Error in querying with rows\n")
		}
	}
	// надо сделать какой нить ID для того чтобы отсканить почеловечески с WHERE в БД для точности
	err = conn.QueryRow(ctx, os.Getenv("GET_USER_ID_QUERY"), email).Scan(&userID)
	if err != nil{
		logger.Err(err).Msg("Error in querying the row in getting user's ID\n")
	}
	return userID
}

func GetNotes() []mod.Note{
	conn, err := ConnectingSQL()
	if err != nil {
		logger.Err(err).Msg("Error in getting notes\n")
	}

	usersData := []mod.Note{} 
	rows, err := conn.Query(ctx, os.Getenv("GET_NOTES"))
	if err != nil{
		logger.Err(err).Msg("Error in querying with rows\n")
	}
	defer rows.Close()

	for rows.Next(){
		note := mod.Note{}
	
		if err := rows.Scan(&note.NotesData, &note.CreatedAt); err!= nil{
			logger.Err(err).Msg("Error in scaning data with rows\n")
		}
		
		usersData = append(usersData, note)
	}
	return usersData
}

// возвращает ID записки
func GetNoteID() (string, uuid.UUID){
	conn, err := ConnectingSQL()
	if err != nil{
		logger.Err(err).Msg("Error in getting ID\n")
	}
	defer conn.Close(ctx)
	
	var (
		noteID string
		noteUUID uuid.UUID
	)
		
	if err := conn.QueryRow(ctx, os.Getenv("GET_NOTE_ID_UUID")).Scan(&noteID, &noteUUID);err != nil{
		logger.Err(err).Msg("Error in querying the row in getting note id/uuid\n")
	}
	return noteID, noteUUID
}

// @Summary User's notes here
// @Description Заметки пользователя
// @Router /users/notes [get]
func ShowNotes(c echo.Context) error{
	info := GetNotes()
	userID := GetUserID()
	user := mod.User{}
	user.ID = userID
	_, UUID := GetNoteID()

	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"Title": "Notes", 
		"Notes": info,
		"UserID": user.ID,
		"NoteID": UUID,
	})
}

// надо через c.Set() в отдельной функции сделать так чтобы получать id из базы данных и вставлять в ссылку users/:id  
// Создать в базе данных ID для заметок - для /users/main/:id - чтобы удалить заметку с этой же id
func DeleteNotes(c echo.Context) error {
	conn, err := ConnectingSQL()
	if err != nil{
		logger.Err(err).Msg("Can't connect to database\n")
	}
	defer conn.Close(ctx)

	_, noteUUID := GetNoteID()
	strUUID := noteUUID.String()
	_, err = conn.Exec(ctx, os.Getenv("DELETE_NOTES"), noteUUID)
	if err != nil{
		logger.Err(err).Msg("Can't delete the note\n")
	}
	return c.Redirect(http.StatusOK, "/users/"+GetUserID()+"/notes/"+strUUID+"/delete")
}

func WriteNotes(c echo.Context) error {
	notes := c.FormValue("write_notes")
	notesUUID := uuid.New()
	strUUID := notesUUID.String()

	if notes != ""{
		conn, err := ConnectingSQL()
		if err != nil{
			logger.Err(err).Msg("Can't connect to database\n")
		}
		defer conn.Close(ctx)

		_, err = conn.Exec(ctx, os.Getenv("WRITE_NOTES"), strUUID, notes)
		if err != nil{
			logger.Err(err).Msg("Can't insert user's notes in DB\n")
		}
		return ShowNotes(c)
	}
	return nil
}