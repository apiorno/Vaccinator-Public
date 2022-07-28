package controllertests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/S-S-Group/Vaccinator/src/models"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var DBClient *gorm.DB

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}
	Database()

	os.Exit(m.Run())

}

func Database() {

	var err error

	TestDbDriver := os.Getenv("TestDbDriver")

	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbUser"), os.Getenv("TestDbName"), os.Getenv("TestDbPassword"))
	DBClient, err = gorm.Open(TestDbDriver, DBURL)
	if err != nil {
		fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the %s database\n", TestDbDriver)
	}

}

func refreshUserTable() error {
	err := DBClient.DropTableIfExists(&models.User{}).Error
	if err != nil {
		return err
	}
	err = DBClient.AutoMigrate(&models.User{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed table")
	return nil
}

func seedOneUser() (models.User, error) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user := models.User{
		Firstname: "Pet",
		Lastname:  "Pit",
		Passport:  "2222",
		Email:     "pet@gmail.com",
		Password:  "password",
	}

	err = DBClient.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func seedUsers() ([]models.User, error) {

	var err error
	if err != nil {
		return nil, err
	}
	users := []models.User{
		{
			Firstname: "Steven victor",
			Lastname:  "Hansen",
			Passport:  "3333",
			Email:     "steven@gmail.com",
			Password:  "password",
		},
		{
			Firstname: "Kenny Morris",
			Lastname:  "Lyl",
			Passport:  "4444",
			Email:     "kenny@gmail.com",
			Password:  "password",
		},
	}
	for i, _ := range users {
		err := DBClient.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			return []models.User{}, err
		}
	}
	return users, nil
}

func refreshUserAndNotificationTable() error {

	err := DBClient.DropTableIfExists(&models.User{}, &models.Notification{}).Error
	if err != nil {
		return err
	}
	err = DBClient.AutoMigrate(&models.User{}, &models.Notification{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed tables")
	return nil
}

func seedOneUserAndOneNotification() (models.Notification, error) {

	err := refreshUserAndNotificationTable()
	if err != nil {
		return models.Notification{}, err
	}
	user := models.User{
		Firstname: "Sam Phil",
		Lastname:  "Callin",
		Passport:  "5555",
		Email:     "sam@gmail.com",
		Password:  "password",
	}
	err = DBClient.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.Notification{}, err
	}
	notification := models.Notification{
		Title:    "This is the title sam",
		Content:  "This is the content sam",
		AuthorID: user.ID,
	}
	err = DBClient.Model(&models.Notification{}).Create(&notification).Error
	if err != nil {
		return models.Notification{}, err
	}
	return notification, nil
}

func seedUsersAndNotifications() ([]models.User, []models.Notification, error) {

	var err error

	if err != nil {
		return []models.User{}, []models.Notification{}, err
	}
	var users = []models.User{
		{
			Firstname: "Steven victor",
			Lastname:  "Hansen",
			Passport:  "6666",
			Email:     "steven@gmail.com",
			Password:  "password",
		},
		{
			Firstname: "Magu Frank",
			Lastname:  "Lyl",
			Passport:  "7777",
			Email:     "magu@gmail.com",
			Password:  "password",
		},
	}
	var notifications = []models.Notification{
		{
			Title:   "Title 1",
			Content: "Hello world 1",
		},
		{
			Title:   "Title 2",
			Content: "Hello world 2",
		},
	}

	for i, _ := range users {
		err = DBClient.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		notifications[i].AuthorID = users[i].ID

		err = DBClient.Model(&models.Notification{}).Create(&notifications[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}
	}
	return users, notifications, nil
}
