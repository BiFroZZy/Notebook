package auth

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"

	mod "notebook/internal/models"
	db "notebook/internal/repository"
)
func HashingFunc(password string) (hashedPass []byte){
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil{
		logger.Err(err).Msg("Error while hashing the password!")
	}
	return hashedPass
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
	
	conn, err := db.ConnectingSQL()
	if err != nil {
		logger.Err(err).Msg("Can't connect to database\n")
	}
	defer conn.Close(ctx)
	
	user := mod.User{}
	err = conn.QueryRow(ctx, os.Getenv("REG_QUERY"), getRegLogin).Scan(&user.Login, &user.Password, &user.Email)
	if err != nil{
		logger.Err(err).Msg("Can't get user's info in registration\n")
	}
	HashPassword := HashingFunc(getRegPassword)
	ok := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(getRegPassword))
	if getRegLogin == user.Login && ok == nil{
		return c.Render(http.StatusOK, "reg.html", map[string]interface{}{
			"Title": "Registration",
			"Error": "Login or password already exist",
		})
	}
	if getRegEmail == user.Email{
		return c.Render(http.StatusOK, "reg.html", map[string]interface{}{
			"Title": "Registration",
			"Error": "You already have an account!",
		})
	}		
	stringUUID := newUUID.String()
	db.WriteDataSQL(stringUUID, getRegLogin, string(HashPassword), getRegEmail)
	return c.Redirect(http.StatusFound, "/users/"+stringUUID+"/notes")
}

