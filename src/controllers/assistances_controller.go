package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/S-S-Group/Vaccinator/src/data"
	"github.com/S-S-Group/Vaccinator/src/models"
	"github.com/gorilla/mux"
	"github.com/S-S-Group/Vaccinator/src/responses"
)

// AssistancesController represents the controller to handle requests for user assistances
type AssistancesController struct {
	l *log.Logger
}

// LogLine logs the text using the controller's logger
func (c *AssistancesController) LogLine(text string) {
	c.l.Println(text)
}

// GetUserById return the user associated to the requested id
func (c *AssistancesController) GetAssistancesOfUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		c.LogLine("Handle GET/{id} ")

		id := mux.Vars(r)["id"]

		assistances, err := data.GetAssistancesOfUser(id)

		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Can not retrieve data from database", http.StatusBadRequest)
			return
		}

		err = json.NewEncoder(rw).Encode(assistances)

		if err != nil {
			http.Error(rw, "Can not convert user to JSON", http.StatusInternalServerError)
			return
		}
	}
}

// GetUserById return the user associated to the requested id
func (c *AssistancesController) GetAssistancesCount() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		c.LogLine("Handle GET/ ")

	

		count, err := data.GetAssistancesCount()

		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Can not retrieve data from database", http.StatusBadRequest)
			return
		}
		responses.JSON(rw, 200, struct {
			Count     int64 `json:"count"`
		}{
			count,
		})
		
	}
}

//CreateAssistance creates user assistance from json data
func (c *AssistancesController) CreateAssistance() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		c.LogLine("Handle POST ")

		assistance := &models.Assistance{}

		err := assistance.FromJSON(r.Body)

		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Unable to parse json", http.StatusBadRequest)
			return
		}
		assistance.Prepare()
		err = assistance.Validate()
		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Invalid certification information", http.StatusBadRequest)
			return
		}
		assistance, err = data.SaveAssistance(assistance)
		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Can not add assistance", http.StatusInternalServerError)
			return
		}
		err = assistance.ToJSON(rw)

		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Can not convert assistance to JSON", http.StatusInternalServerError)
			return
		}
	}
}

// DeleteAssistance removes the assistance associated with  the requested id
func (c *AssistancesController) DeleteAssistance() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		c.LogLine("Handle DELETE/{id} ")

		id := mux.Vars(r)["id"]

		_, err := data.DeleteAssistance(id)
		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Can not remove assistance", http.StatusBadRequest)
			return
		}
		rw.Header().Set("Entity", fmt.Sprintf("%d", id))

	}
}

// UpdateAssistance updates the assistance associated with the requested id
func (c *AssistancesController) UpdateAssistance() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		c.LogLine("Handle UPDATE/{id} ")

		id := mux.Vars(r)["id"]

		assistance := &models.Assistance{}
		err := assistance.FromJSON(r.Body)

		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Unable to parse json", http.StatusBadRequest)
			return
		}

		assistance.Prepare()
		err = assistance.Validate()
		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Invalid assistance information", http.StatusBadRequest)
			return
		}

		updatedAssistance, err := data.UpdateAssistance(id, assistance)

		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Can not update assistance", http.StatusInternalServerError)
			return
		}
		err = updatedAssistance.ToJSON(rw)

		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Can not convert assistance to JSON", http.StatusInternalServerError)
			return
		}
	}
}
