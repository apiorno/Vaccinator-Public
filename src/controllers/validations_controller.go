package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/S-S-Group/Vaccinator/src/data"
	"github.com/S-S-Group/Vaccinator/src/responses"
)

// CertificationsController represents the controller to handle requests for URLs
type ValidationsController struct {
	l *log.Logger
}
type UserValidationBody struct {
	UserID string `json:"userID"`
}

// LogLine logs the text using the controller's logger
func (c *ValidationsController) LogLine(text string) {
	c.l.Println(text)
}

// GetUserById return the user associated to the requested id
func (c *ValidationsController) ValidateUserCertifications() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		c.LogLine("Handle POST/ ")
		userValidation := &UserValidationBody{}
		err := json.NewDecoder(r.Body).Decode(userValidation)

		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Malformed body", http.StatusBadRequest)
			return
		}

		stringId := userValidation.UserID

		c.LogLine("INICIO CONSULTA CERT")
		certification, err := data.GetCertificationsOfUser(stringId)
		c.LogLine("FIN CONSULTA CERT")
		state := "PENDING"
		if (err != nil) && (err.Error() == "sql: no rows in result set"){
			state = "UNREGISTERED"
		}
		if (err != nil) && (err.Error() != "sql: no rows in result set") {
			
			c.LogLine(err.Error())
			http.Error(rw, "Can not retrieve data from database", http.StatusBadRequest)
			return
		}

	        c.LogLine("INICIO CONSULTA USER")
		user, err := data.FindUserById(stringId, 0)
		c.LogLine("FIN CONSULTAVUSER")
		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Can not retrieve data from database", http.StatusBadRequest)
			return
		}
		
		if certification!= nil && certification.Validated {
			state = "VALIDATED"
		}
		
		
		
		c.LogLine("FIN DEL FIN")
		responses.JSON(rw, 200, struct {
			Valid     string `json:"state"`
			Firstname string `json:"firstname"`
			Lastname  string `json:"lastname"`
			Passport  string `json:"passport"`
		}{
			state,
			user.Firstname,
			user.Lastname,
			user.Passport,
		})

	}
}
