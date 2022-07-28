package modeltests

import (
	"log"
	"testing"

	"github.com/S-S-Group/Vaccinator/src/data"
	"github.com/S-S-Group/Vaccinator/src/models"
	"gopkg.in/go-playground/assert.v1"
)

func TestFindAllUsers(t *testing.T) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	err = seedUsers()
	if err != nil {
		log.Fatal(err)
	}

	users, err := data.FindAllUsers()
	if err != nil {
		t.Errorf("this is the error getting the users: %v\n", err)
		return
	}
	assert.Equal(t, len(*users), 2)
}

func TestSaveUser(t *testing.T) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	newUser := models.User{
		ID:        1,
		Email:     "test@gmail.com",
		Firstname: "test",
		Lastname:  "test",
		Passport:  "1",
		Password:  "password",
	}
	savedUser, err := data.SaveUser(&newUser)
	if err != nil {
		t.Errorf("this is the error getting the users: %v\n", err)
		return
	}
	assert.Equal(t, newUser.ID, savedUser.ID)
	assert.Equal(t, newUser.Email, savedUser.Email)
	assert.Equal(t, newUser.Firstname, savedUser.Firstname)
	assert.Equal(t, newUser.Lastname, savedUser.Lastname)
	assert.Equal(t, newUser.Passport, savedUser.Passport)
}

func TestGetUserByID(t *testing.T) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("cannot seed users table: %v", err)
	}
	foundUser, err := data.FindUserById(user.ID, 1)
	if err != nil {
		t.Errorf("this is the error getting one user: %v\n", err)
		return
	}
	assert.Equal(t, foundUser.ID, user.ID)
	assert.Equal(t, foundUser.Email, user.Email)
	assert.Equal(t, foundUser.Firstname, user.Firstname)
	assert.Equal(t, foundUser.Lastname, user.Lastname)
	assert.Equal(t, foundUser.Passport, user.Passport)
}

func TestUpdateAUser(t *testing.T) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user: %v\n", err)
	}

	userUpdate := models.User{
		ID:        1,
		Firstname: "modiUpdate",
		Lastname:  "modiUpdate",
		Passport:  "1",
		Email:     "modiupdate@gmail.com",
		Password:  "password",
	}
	updatedUser, err := data.UpdateUser(user.ID, &userUpdate)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	assert.Equal(t, updatedUser.ID, userUpdate.ID)
	assert.Equal(t, updatedUser.Email, userUpdate.Email)
	assert.Equal(t, updatedUser.Firstname, userUpdate.Firstname)
	assert.Equal(t, updatedUser.Lastname, userUpdate.Lastname)
	assert.Equal(t, updatedUser.Passport, userUpdate.Passport)
}

func TestDeleteAUser(t *testing.T) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()

	if err != nil {
		log.Fatalf("Cannot seed user: %v\n", err)
	}

	isDeleted, err := data.DeleteUser(user.ID)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}

	assert.Equal(t, isDeleted, int64(1))
}
