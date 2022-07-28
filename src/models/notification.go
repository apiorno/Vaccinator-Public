package models

import (
	"encoding/json"
	"errors"
	"html"
	"io"
	"strings"
	"time"
)

type Notification struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Title     string    `gorm:"size:255;not null;unique" json:"title"`
	Content   string    `gorm:"size:255;not null;" json:"content"`
	Author    User      `json:"author"`
	AuthorID  uint32    `gorm:"not null" json:"author_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (n *Notification) Prepare() {
	n.ID = 0
	n.Title = html.EscapeString(strings.TrimSpace(n.Title))
	n.Content = html.EscapeString(strings.TrimSpace(n.Content))
	n.Author = User{}
	n.CreatedAt = time.Now()
	n.UpdatedAt = time.Now()
}

func (n *Notification) Validate() error {

	if n.Title == "" {
		return errors.New("Required Title")
	}
	if n.Content == "" {
		return errors.New("Required Content")
	}
	if n.AuthorID < 1 {
		return errors.New("Required Author")
	}
	return nil
}

// ToJSON generates a json representation of Notification
func (u *Notification) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(u)
}

// FromJSON generates a Notification from a json
func (u *Notification) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(u)
}
