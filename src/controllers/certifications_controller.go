package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/S-S-Group/Vaccinator/src/data"
	"github.com/S-S-Group/Vaccinator/src/models"
	"github.com/gorilla/mux"
)

// CertificationsController represents the controller to handle requests for URLs
type CertificationsController struct {
	l *log.Logger
}

// LogLine logs the text using the controller's logger
func (c *CertificationsController) LogLine(text string) {
	c.l.Println(text)
}

// GetUserById return the user associated to the requested id
func (c *CertificationsController) GetCertificationsOfUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		c.LogLine("Handle GET/{id} ")

		id := mux.Vars(r)["id"]
		c.LogLine("INICIO FOTO")
		certification, err := data.GetCertificationsOfUser(id)
		c.LogLine("FIN FOTO")
		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Can not retrieve data from database", http.StatusBadRequest)
			return
		}

		err = json.NewEncoder(rw).Encode(certification)

		if err != nil {
			http.Error(rw, "Can not convert user to JSON", http.StatusInternalServerError)
			return
		}
	}
}

//CreateUser creates user from json data
func (c *CertificationsController) CreateCertification() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		c.LogLine("Handle POST ")

		certification := &models.Certification{}

		err := certification.FromJSON(r.Body)
		c.LogLine("INICIO")
		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Unable to parse json", http.StatusBadRequest)
			return
		}
		c.LogLine("PREPARE")
		certification.Prepare()
		err = certification.Validate()
		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Invalid certification information", http.StatusBadRequest)
			return
		}
		c.LogLine("ENTRANDO")
		certification, err = data.SaveCertification(certification)
		c.LogLine("SALIENDO")
		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Can not add certification", http.StatusInternalServerError)
			return
		}
		err = certification.ToJSON(rw)

		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Can not convert certification to JSON", http.StatusInternalServerError)
			return
		}
	}
}

// DeleteUser removes the user associated with  the requested id
func (c *CertificationsController) DeleteCertification() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		c.LogLine("Handle DELETE/{id} ")

		id := mux.Vars(r)["id"]

		_, err := data.DeleteCertification(id)
		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Can not remove user", http.StatusBadRequest)
			return
		}
		rw.Header().Set("Entity", fmt.Sprintf("%d", id))

	}
}

// UpdateUser updates the user associated with the requested id
func (c *CertificationsController) UpdateCertification() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		c.LogLine("Handle UPDATE/{id} ")

		id := mux.Vars(r)["id"]

		certification := &models.Certification{}
		err := certification.FromJSON(r.Body)

		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Unable to parse json", http.StatusBadRequest)
			return
		}

		certification.Prepare()
		err = certification.Validate()
		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Invalid certification information", http.StatusBadRequest)
			return
		}

		updatedCertification, err := data.UpdateCertification(id, certification)

		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Can not update certification", http.StatusInternalServerError)
			return
		}
		err = updatedCertification.ToJSON(rw)

		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Can not convert certification to JSON", http.StatusInternalServerError)
			return
		}
	}
}
