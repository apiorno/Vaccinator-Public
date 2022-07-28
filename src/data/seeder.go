package data

import (
	"log"

	"github.com/S-S-Group/Vaccinator/src/models"
	"github.com/jinzhu/gorm"
)

var users = []models.User{
	{
		Firstname: "Steven Victor",
		Lastname:  "Hansen",
		Passport:  "202020",
		Email:     "steven@gmail.com",
		Password:  "password",
		Role:      1,
	},
	{
		Firstname: "Martin Luther",
		Lastname:  "King",
		Passport:  "303030",
		Email:     "luther@gmail.com",
		Password:  "password",
		Role:      2,
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

var certifications = []models.Certification{
	{
		UserID:    "2345",
		Bytes:     "asdasdas",
		Validated: true,
	},
	{
		UserID:    "bgf1246",
		Bytes:     "werrdfg",
		Validated: false,
	},
}

func Load(db *gorm.DB, l *log.Logger) {

	/*err := db.Debug().DropTableIfExists(&models.Notification{}, &models.User{}, &models.Certification{}, &models.Assistance{}).Error
	if err != nil {
		l.Fatalf("cannot drop table: %v", err)
	}*/
	err := db.Debug().AutoMigrate(&models.Notification{}, &models.User{}, &models.Certification{}, &models.Assistance{}).Error
	if err != nil {
		l.Fatalf("cannot migrate table: %v", err)
	}
	err = db.Debug().Model(&models.Notification{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	/*for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			l.Fatalf("cannot seed users table: %v", err)
		}
		notifications[i].AuthorID = users[i].ID

		err = db.Debug().Model(&models.Notification{}).Create(&notifications[i]).Error
		if err != nil {
			log.Fatalf("cannot seed notifications table: %v", err)
		}
	}
	for i, _ := range certifications {
		err = db.Debug().Model(&models.Certification{}).Create(&certifications[i]).Error
		if err != nil {
			l.Fatalf("cannot seed certifications table: %v", err)
		}
	}*/
}
