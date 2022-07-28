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
type NotificationsController struct {
	l *log.Logger
}

// LogLine logs the text using the controller's logger
func (c *NotificationsController) LogLine(text string) {
	c.l.Println(text)
}

// GetAllUsers return all the users
func (c *NotificationsController) GetAllNotifications() http.HandlerFunc {
	c.LogLine("GET ALL USERS")
	return func(rw http.ResponseWriter, r *http.Request) {

		c.LogLine("Handle GET")

		users, err := data.FindAllNotifications()
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
func (c *NotificationsController) GetNotificationById() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		c.LogLine("Handle GET/{id} ")

		stringId := mux.Vars(r)["id"]
		id, err := strconv.Atoi(stringId)

		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Id should be a number", http.StatusBadRequest)
			return
		}

		notification, err := data.FindNotificationById(uint32(id))

		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Can not retrieve data from database", http.StatusBadRequest)
			return
		}

		err = json.NewEncoder(rw).Encode(notification)

		if err != nil {
			http.Error(rw, "Can not convert user to JSON", http.StatusInternalServerError)
			return
		}
	}
}

//CreateUser creates user from json data
func (c *NotificationsController) CreateNotification() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		c.LogLine("Handle POST ")

		notification := &models.Notification{}
		err := notification.FromJSON(r.Body)

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

		if tokenID != uint32(notification.AuthorID) {
			c.LogLine(err.Error())
			http.Error(rw, "Unauthorized", http.StatusUnauthorized)
			return
		}
		notification.Prepare()
		err = notification.Validate()
		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Invalid user information", http.StatusBadRequest)
			return
		}
		notification, err = data.SaveNotification(notification)
		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Can not add user", http.StatusInternalServerError)
			return
		}
		rw.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, notification.ID))
		err = notification.ToJSON(rw)

		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Can not convert user to JSON", http.StatusInternalServerError)
			return
		}
	}
}

// UpdateUser updates the user associated with the requested id
func (c *NotificationsController) UpdateNotification() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		c.LogLine("Handle UPDATE/{id} ")

		stringId := mux.Vars(r)["id"]
		id, err := strconv.Atoi(stringId)

		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Invalid id type", http.StatusBadRequest)
			return
		}

		notification := &models.Notification{}
		err = notification.FromJSON(r.Body)

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
		notification.Prepare()
		err = notification.Validate()
		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Invalid user information", http.StatusBadRequest)
			return
		}

		updatedNotification, err := data.UpdateNotification(uint32(id), notification)

		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Can not update user", http.StatusInternalServerError)
			return
		}
		err = updatedNotification.ToJSON(rw)

		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Can not convert notification to JSON", http.StatusInternalServerError)
			return
		}
	}
}

// DeleteUser removes the user associated with  the requested id
func (c *NotificationsController) DeleteNotification() http.HandlerFunc {
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
		_, err = data.DeleteNotification(uint32(id))
		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Can not remove notification", http.StatusBadRequest)
			return
		}
		rw.Header().Set("Entity", fmt.Sprintf("%d", id))

	}
}
