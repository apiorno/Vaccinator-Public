package modeltests

import (
	"log"
	"testing"

	"github.com/S-S-Group/Vaccinator/src/data"
	"github.com/S-S-Group/Vaccinator/src/models"
	"gopkg.in/go-playground/assert.v1"
)

func TestFindAllNotifications(t *testing.T) {

	err := refreshUserAndNotificationTable()
	if err != nil {
		log.Fatalf("Error refreshing user and notification table %v\n", err)
	}
	_, _, err = seedUsersAndNotifications()
	if err != nil {
		log.Fatalf("Error seeding user and notification  table %v\n", err)
	}
	notifications, err := data.FindAllNotifications()
	if err != nil {
		t.Errorf("this is the error getting the notifications: %v\n", err)
		return
	}
	assert.Equal(t, len(*notifications), 2)
}

func TestSavePost(t *testing.T) {

	err := refreshUserAndNotificationTable()
	if err != nil {
		log.Fatalf("Error user and notification refreshing table %v\n", err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}

	newNotification := models.Notification{
		ID:       1,
		Title:    "This is the title",
		Content:  "This is the content",
		AuthorID: user.ID,
	}
	savedNotification, err := data.SaveNotification(&newNotification)
	if err != nil {
		t.Errorf("this is the error getting the post: %v\n", err)
		return
	}
	assert.Equal(t, newNotification.ID, savedNotification.ID)
	assert.Equal(t, newNotification.Title, savedNotification.Title)
	assert.Equal(t, newNotification.Content, savedNotification.Content)
	assert.Equal(t, newNotification.AuthorID, savedNotification.AuthorID)

}

func TestGetPostByID(t *testing.T) {

	err := refreshUserAndNotificationTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table: %v\n", err)
	}
	notification, err := seedOneUserAndOneNotification()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	foundNotification, err := data.FindNotificationById(uint32(notification.ID))
	if err != nil {
		t.Errorf("this is the error getting one user: %v\n", err)
		return
	}
	assert.Equal(t, foundNotification.ID, notification.ID)
	assert.Equal(t, foundNotification.Title, notification.Title)
	assert.Equal(t, foundNotification.Content, notification.Content)
}

func TestUpdateAPost(t *testing.T) {

	err := refreshUserAndNotificationTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table: %v\n", err)
	}
	notification, err := seedOneUserAndOneNotification()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	notificationUpdate := models.Notification{
		ID:       1,
		Title:    "modiUpdate",
		Content:  "modiupdate@gmail.com",
		AuthorID: notification.AuthorID,
	}
	updatedNotification, err := data.UpdateNotification(uint32(notification.ID), &notificationUpdate)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	assert.Equal(t, updatedNotification.ID, notification.ID)
	assert.Equal(t, updatedNotification.Title, notification.Title)
	assert.Equal(t, updatedNotification.Content, notification.Content)
	assert.Equal(t, updatedNotification.AuthorID, notification.AuthorID)
}

func TestDeleteAPost(t *testing.T) {

	err := refreshUserAndNotificationTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table: %v\n", err)
	}
	notification, err := seedOneUserAndOneNotification()
	if err != nil {
		log.Fatalf("Error Seeding tables")
	}
	isDeleted, err := data.DeleteNotification(uint32(notification.ID))
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}

	assert.Equal(t, isDeleted, int64(1))
}
