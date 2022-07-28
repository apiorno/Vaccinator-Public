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
type Certification struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"-"`
	UserID    string    `gorm:"size:255;not null" json:"userID"`
	Bytes     string    `gorm:"not null" json:"bytes"`
	Date      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"date"`
	Validated bool      `gorm:"default:FALSE" json:"validated"`
}

func (c *Certification) Prepare() {
	c.ID = 0
	c.UserID = html.EscapeString(strings.TrimSpace(c.UserID))
	c.Date = time.Now()
}

func (c *Certification) Validate() error {

	if c.UserID == "" {
		return errors.New("Required UserID")
	}
	if c.Bytes == "" {
		return errors.New("Required image")
	}

	return nil

}

// ToJSON generates a json representation of User
func (c *Certification) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(c)
}

// FromJSON generates an User from a json
func (c *Certification) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(c)
}
