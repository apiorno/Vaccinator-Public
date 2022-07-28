package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/S-S-Group/Vaccinator/src/auth"
	"github.com/S-S-Group/Vaccinator/src/data"
	"github.com/S-S-Group/Vaccinator/src/models"
	"github.com/gorilla/mux"
)

// UsersController represents the controller to handle requests for URLs
type UsersController struct {
	l *log.Logger
}

// LogLine logs the text using the controller's logger
func (c *UsersController) LogLine(text string) {
	c.l.Println(text)
}

// GetAllUsers return all the users
func (c *UsersController) GetAllUsers() http.HandlerFunc {
	c.LogLine("GET ALL USERS")
	return func(rw http.ResponseWriter, r *http.Request) {

		c.LogLine("Handle GET")

		users, err := data.FindAllUsers()
		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Can not retrieve data from database", http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(rw).Encode(users)

		if err != nil {
			http.Error(rw, "Can not convert users to JSON", http.StatusInternalServerError)
			return
		}
	}
}

// GetUserById return the user associated to the requested id
func (c *UsersController) GetUserById() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		c.LogLine("Handle GET/{id} ")
		//role, err := strconv.Atoi(r.URL.Query()["role"][0])
		/*if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Role should be a number", http.StatusBadRequest)
			return
		}*/
		id := mux.Vars(r)["id"]

		user, err := data.FindUserById(id, 0)

		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Can not retrieve data from database", http.StatusBadRequest)
			return
		}

		err = json.NewEncoder(rw).Encode(user)

		if err != nil {
			http.Error(rw, "Can not convert user to JSON", http.StatusInternalServerError)
			return
		}
	}
}

//CreateUser creates user from json data
func (c *UsersController) CreateUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		c.LogLine("Handle POST ")

		user := &models.User{}
		err := user.FromJSON(r.Body)

		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Unable to parse json", http.StatusBadRequest)
			return
		}
		user.Prepare()
		err = user.Validate()
		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Invalid user information", http.StatusBadRequest)
			return
		}
		user, err = data.SaveUser(user)
		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Can not add user", http.StatusInternalServerError)
			return
		}
		rw.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, user.ID))
		err = user.ToJSON(rw)

		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Can not convert user to JSON", http.StatusInternalServerError)
			return
		}
	}
}

// DeleteUser removes the user associated with  the requested id
func (c *UsersController) DeleteUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		c.LogLine("Handle DELETE/{id} ")

		stringId := mux.Vars(r)["id"]
		id, err := strconv.Atoi(stringId)

		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Invalid id type", http.StatusBadRequest)
			return
		}

		tokenID, err := auth.ExtractTokenID(r)
		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if tokenID != 0 && tokenID != uint32(id) {
			c.LogLine(err.Error())
			http.Error(rw, "Unauthorized", http.StatusUnauthorized)
			return
		}
		_, err = data.DeleteUser(stringId)
		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Can not remove user", http.StatusBadRequest)
			return
		}
		rw.Header().Set("Entity", fmt.Sprintf("%d", id))

	}
}

// UpdateUser updates the user associated with the requested id
func (c *UsersController) UpdateUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		c.LogLine("Handle UPDATE/{id} ")

		stringId := mux.Vars(r)["id"]
		id, err := strconv.Atoi(stringId)

		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Invalid id type", http.StatusBadRequest)
			return
		}

		user := &models.User{}
		err = user.FromJSON(r.Body)

		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Unable to parse json", http.StatusBadRequest)
			return
		}

		tokenID, err := auth.ExtractTokenID(r)
		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if tokenID != uint32(id) {
			c.LogLine(err.Error())
			http.Error(rw, "Unauthorized", http.StatusUnauthorized)
			return
		}
		user.Prepare()
		err = user.Validate()
		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Invalid user information", http.StatusBadRequest)
			return
		}

		updatedUser, err := data.UpdateUser(stringId, user)

		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Can not update user", http.StatusInternalServerError)
			return
		}
		err = updatedUser.ToJSON(rw)

		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Can not convert user to JSON", http.StatusInternalServerError)
			return
		}
	}
}
