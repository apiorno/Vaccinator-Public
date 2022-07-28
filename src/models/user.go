package models

import (
	"encoding/json"
	"errors"
	"html"
	"io"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"golang.org/x/crypto/bcrypt"
)

// User represents the user attending an event
type User struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"-"`
	Firstname string    `gorm:"size:255;not null" json:"firstname"`
	Lastname  string    `gorm:"size:255;not null" json:"lastname"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Passport  string    `gorm:"size:50;not null;unique" json:"passport"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	Role      int       `gorm:"default:0" json:"role"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
}

// Hash encrypt password
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// VerifyPassword verfies if the given password is the same as the encrypted one
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *User) BeforeSave() error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) Prepare() {
	u.ID = 0
	u.Firstname = html.EscapeString(strings.TrimSpace(u.Firstname))
	u.Lastname = html.EscapeString(strings.TrimSpace(u.Lastname))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.Passport = html.EscapeString(strings.TrimSpace(u.Passport))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) ValidateLogin() error {

	if u.Password == "" {
		return errors.New("Required Password")
	}
	if u.Email == "" {
		return errors.New("Required Email")
	}
	if err := checkmail.ValidateFormat(u.Email); err != nil {
		return errors.New("Invalid Email")
	}

	return nil

}

func (u *User) Validate() error {

	if u.Firstname == "" {
		return errors.New("Required Firstname")
	}
	if u.Lastname == "" {
		return errors.New("Required Lastname")
	}
	if u.Passport == "" {
		return errors.New("Required Passport")
	}
	if u.Password == "" {
		return errors.New("Required Password")
	}
	if u.Email == "" {
		return errors.New("Required Email")
	}
	if err := checkmail.ValidateFormat(u.Email); err != nil {
		return errors.New("Invalid Email")
	}

	return nil

}

// ToJSON generates a json representation of User
func (u *User) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(u)
}

// FromJSON generates an User from a json
func (u *User) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(u)
}
