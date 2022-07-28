package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/S-S-Group/Vaccinator/src/data"
	"github.com/S-S-Group/Vaccinator/src/models"
	"golang.org/x/crypto/bcrypt"
)

// LoginController represents the controller to handle requests for URLs
type LoginController struct {
	l *log.Logger
}

// LogLine logs the text using the controller's logger
func (c *LoginController) LogLine(text string) {
	c.l.Println(text)
}

func (c *LoginController) Login(rw http.ResponseWriter, r *http.Request) {
	c.LogLine("Handle LOGIN")

	user := &models.User{}
	err := user.FromJSON(r.Body)

	if err != nil {
		c.LogLine(err.Error())
		http.Error(rw, "Unable to parse json", http.StatusBadRequest)
		return
	}

	user.Prepare()
	err = user.ValidateLogin()
	if err != nil {
		c.LogLine(err.Error())
		http.Error(rw, "Invalid user information", http.StatusBadRequest)
		return
	}
	resultUser, err := c.SignIn(user.Email, user.Password)
	if err != nil {
		c.LogLine(err.Error())
		http.Error(rw, "Failed sign in", http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(rw).Encode(resultUser)

	if err != nil {
		http.Error(rw, "Can not convert token to JSON", http.StatusInternalServerError)
		return
	}
}

func (c *LoginController) SignIn(email, password string) (*models.User, error) {

	var err error

	user := models.User{}

	userSql := "SELECT  email, password, passport,firstname,lastname FROM users WHERE email = $1"

	err = data.DBClient2.QueryRow(userSql, email).Scan(&user.Email, &user.Password, &user.Passport,&user.Firstname,&user.Lastname)
	if err != nil {
		return nil, err
	}

	/*err = data.DBClient.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return nil, err
	}*/
	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, err
	}
	//return auth.CreateToken(user.ID)
	return &user, nil
}
