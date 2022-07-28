package models

import (
	"encoding/json"
	"errors"
	"html"
	"io"
	"strings"
	"time"
)

// User represents the user attending an event
type Assistance struct {
	ID           uint32    `gorm:"primary_key;auto_increment" json:"-"`
	UserID       string    `gorm:"size:255;not null" json:"userID"`
	SupervisorID string    `gorm:"size:255;not null" json:"supervisorID"`
	Date         time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"date"`
}

func (a *Assistance) Prepare() {
	a.ID = 0
	a.UserID = html.EscapeString(strings.TrimSpace(a.UserID))
	a.SupervisorID = html.EscapeString(strings.TrimSpace(a.SupervisorID))
	a.Date = time.Now()
}

func (a *Assistance) Validate() error {

	if a.UserID == "" {
		return errors.New("Required UserID")
	}
	if a.SupervisorID == "" {
		return errors.New("Required SupervisorID")
	}

	return nil

}

// ToJSON generates a json representation of User
func (a *Assistance) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(a)
}

// FromJSON generates an User from a json
func (a *Assistance) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(a)
}
