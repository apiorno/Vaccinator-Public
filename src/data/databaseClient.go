package data

import (
	"fmt"
	"log"
	"time"
	"github.com/S-S-Group/Vaccinator/src/models"
	"github.com/jinzhu/gorm"
	"database/sql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
)

// ErrorUserNotFound is the common error message when a user is not found
var ErrorUserNotFound = fmt.Errorf("User not found")

// ErrorNOtificationNotFound is the common error message when a notification is not found
var ErrorNOtificationNotFound = fmt.Errorf("Notification not found")
var DBClient *gorm.DB
var DBClient2 *sql.DB

func Connect2(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string, l *log.Logger) *sql.DB {
	var err error

	DBURL := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", DbUser, DbPassword, DbHost, DbPort, DbName)
	DBClient2, err = sql.Open(Dbdriver, DBURL)
	if err != nil {
		l.Printf("Cannot connect to %s database because %v", Dbdriver, err)
		return nil
	}
	l.Printf("We are connected to the %s database", Dbdriver)

	return DBClient2

}

// Connect tries to connect to postgre bd using Gorm
func Connect(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string, l *log.Logger) *gorm.DB {
	var err error

	DBURL := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", DbUser, DbPassword, DbHost, DbPort, DbName)
	DBClient, err = gorm.Open(Dbdriver, DBURL)
	if err != nil {
		l.Printf("Cannot connect to %s database because %v", Dbdriver, err)
		return nil
	}
	l.Printf("We are connected to the %s database", Dbdriver)
	DBClient.Debug().AutoMigrate(&models.User{}, &models.Notification{}, &models.Certification{}) //database migration
	return DBClient

}

// GetAssistancesOfUser return all the assistances for User
func GetAssistancesOfUser(userID string) (*[]models.Assistance, error) {

	var err error
	assistances := []models.Assistance{}
	err = DBClient.Debug().Model(&models.Assistance{}).Limit(100).Where("user_id = ?", userID).Find(&assistances).Error
	if err != nil {
		return &[]models.Assistance{}, err
	}
	return &assistances, err

}

// GetCertificationsOfUser return all the certifications for User
func GetAssistancesCount() (int64, error) {

	var err error
	var count int64

	userSql := "SELECT  COUNT(*) FROM assistances "

	err = DBClient2.QueryRow(userSql).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, err

}

// SaveAssistances adds the certification
func SaveAssistance(assistance *models.Assistance) (*models.Assistance, error) {

	sqlStatement := `
	INSERT INTO assistances(user_id, supervisor_id, date)
	VALUES ($1, $2, $3)`

_, err := DBClient2.Exec(sqlStatement,assistance.UserID,assistance.SupervisorID,assistance.Date)

	/*err = DBClient.Debug().Create(&user).Error*/
	if err != nil {
		return nil, err
	}
	return assistance, nil



	/*var err error
	err = DBClient.Debug().Create(&assistance).Error
	if err != nil {
		return nil, err
	}
	return assistance, nil*/
}

func DeleteAssistance(uid string) (int64, error) {

	DBClient = DBClient.Debug().Model(&models.Assistance{}).Where("id = ?", uid).Take(&models.Assistance{}).Delete(&models.Assistance{})

	if DBClient.Error != nil {
		if gorm.IsRecordNotFoundError(DBClient.Error) {
			return 0, ErrorNOtificationNotFound
		}
		return 0, DBClient.Error
	}
	return DBClient.RowsAffected, nil
}

func UpdateAssistance(uid string, assistance *models.Assistance) (*models.Assistance, error) {

	DBClient = DBClient.Debug().Model(&models.Assistance{}).Where("id = ?", uid).Updates(assistance)
	if DBClient.Error != nil {
		return nil, DBClient.Error
	}

	return assistance, nil
}

// GetCertificationsOfUser return all the certifications for User
func GetCertificationsOfUser(userID string) (*models.Certification, error) {

	var err error
	certification := models.Certification{}

	userSql := "SELECT  user_id, bytes, validated FROM certifications WHERE user_id = $1"

	err = DBClient2.QueryRow(userSql, userID).Scan(&certification.UserID, &certification.Bytes, &certification.Validated)
	if err != nil {
		return nil, err
	}


	/*err = DBClient.Debug().Model(&models.Certification{}).Limit(100).Where("user_id = ?", userID).Find(&certifications).Error
	if err != nil {
		return &[]models.Certification{}, err
	}*/



	return &certification, err

}

// SaveCertification adds the certification
func SaveCertification(certification *models.Certification) (*models.Certification, error) {

	var err error
	sqlStatement := `
	INSERT INTO certifications(user_id, bytes, date,validated)
	VALUES ($1, $2, $3, $4)`

_, err = DBClient2.Exec(sqlStatement,certification.UserID,certification.Bytes,certification.Date,certification.Validated)

	/*err = DBClient.Debug().Create(&user).Error*/
	if err != nil {
		return nil, err
	}
	return certification, nil



	/*var err error
	err = DBClient.Debug().Create(&certification).Error
	if err != nil {
		return nil, err
	}
	return certification, nil*/
}

func DeleteCertification(uid string) (int64, error) {

	DBClient = DBClient.Debug().Model(&models.Certification{}).Where("id = ?", uid).Take(&models.Certification{}).Delete(&models.Certification{})

	if DBClient.Error != nil {
		if gorm.IsRecordNotFoundError(DBClient.Error) {
			return 0, ErrorNOtificationNotFound
		}
		return 0, DBClient.Error
	}
	return DBClient.RowsAffected, nil
}

func UpdateCertification(uid string, certification *models.Certification) (*models.Certification, error) {
	sqlStatement := `
UPDATE certifications
SET validated = $2
WHERE user_id = $1;`
_, err := DBClient2.Exec(sqlStatement, uid, certification.Validated)

	/*DBClient = DBClient.Debug().Model(&models.Certification{}).Where("user_id = ?", uid).Update("validated", certification.Validated)
	if DBClient.Error != nil {
		return nil, DBClient.Error
	}*/
	if err != nil{
		return nil,err
	}

	return certification, nil
}

// FindAllNotifications return all the users
func FindAllNotifications() (*[]models.Notification, error) {

	var err error
	notifications := []models.Notification{}
	err = DBClient.Debug().Model(&models.Notification{}).Limit(100).Find(&notifications).Error
	if err != nil {
		return &[]models.Notification{}, err
	}
	return &notifications, err

}

// FindUserById finds a user by id and returns error if not found
func FindNotificationById(uid uint32) (*models.Notification, error) {
	var err error
	notification := models.Notification{}
	err = DBClient.Debug().Model(models.Notification{}).Where("id = ?", uid).Take(&notification).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, ErrorNOtificationNotFound
	}
	if notification.ID != 0 {
		err = DBClient.Debug().Model(&models.Notification{}).Where("id = ?", notification.AuthorID).Take(&notification.Author).Error
		if err != nil {
			return &models.Notification{}, err
		}
	}
	return &notification, nil
}

// SaveNotification adds the notification
func SaveNotification(notification *models.Notification) (*models.Notification, error) {
	var err error
	err = DBClient.Debug().Create(&notification).Error
	if err != nil {
		return nil, err
	}
	if notification.ID != 0 {
		err = DBClient.Debug().Model(&models.User{}).Where("id = ?", notification.AuthorID).Take(&notification.Author).Error
		if err != nil {
			return &models.Notification{}, err
		}
	}
	return notification, nil
}

// DeleteUser removes the user associated to the requestd id if exists
func DeleteNotification(uid uint32) (int64, error) {

	DBClient = DBClient.Debug().Model(&models.Notification{}).Where("id = ?", uid).Take(&models.Notification{}).Delete(&models.Notification{})

	if DBClient.Error != nil {
		if gorm.IsRecordNotFoundError(DBClient.Error) {
			return 0, ErrorNOtificationNotFound
		}
		return 0, DBClient.Error
	}
	return DBClient.RowsAffected, nil
}

// UpdateUser updates the user associated to the requested id
func UpdateNotification(uid uint32, notification *models.Notification) (*models.Notification, error) {

	var err error
	DBClient = DBClient.Debug().Model(&models.Notification{}).Where("id = ?", uid).Updates(notification)
	if DBClient.Error != nil {
		return nil, DBClient.Error
	}
	if notification.ID != 0 {
		err = DBClient.Debug().Model(&models.Notification{}).Where("id = ?", notification.AuthorID).Take(&notification.Author).Error
		if err != nil {
			return nil, err
		}
	}
	return notification, nil
}

// FindAllUsers return all the users
func FindAllUsers() (*[]models.User, error) {

	var err error
	users := []models.User{}
	err = DBClient.Debug().Model(&models.User{}).Limit(100).Find(&users).Error
	if err != nil {
		return &[]models.User{}, err
	}
	return &users, err

}

// FindUserById finds a user by id and returns error if not found
func FindUserById(uid string, role int) (*models.User, error) {
	var err error
	user := models.User{}
	//err = DBClient.Debug().Model(models.User{}).Where("passport = ? and role = ?", uid, role).Take(&user).Error
	/*err =DBClient.Raw("SELECT * FROM users WHERE passport = ?", uid ).Scan(&user).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, fmt.Errorf("User Not Found")
	}*/

	userSql := "SELECT  email, password, passport,firstname,lastname FROM users WHERE passport = $1"

	err = DBClient2.QueryRow(userSql, uid).Scan(&user.Email, &user.Password, &user.Passport,&user.Firstname,&user.Lastname)
	if err != nil {
		return nil, err
	}
	return &user, err
}

// SaveUser adds the user
func SaveUser(user *models.User) (*models.User, error) {
	var err error
	sqlStatement := `
INSERT INTO users (email, password, passport,firstname,lastname)
VALUES ($1,$2,$3,$4,$5)`
_, err = DBClient2.Exec(sqlStatement,user.Email, user.Password, user.Passport, user.Firstname,user.Lastname)

	/*err = DBClient.Debug().Create(&user).Error*/
	if err != nil {
		return nil, err
	}
	return user, nil
}

// DeleteUser removes the user associated to the requestd id if exists
func DeleteUser(uid string) (int64, error) {

	DBClient = DBClient.Debug().Model(&models.User{}).Where("id = ?", uid).Take(&models.User{}).Delete(&models.User{})

	if DBClient.Error != nil {
		return 0, DBClient.Error
	}
	return DBClient.RowsAffected, nil
}

// UpdateUser updates the user associated to the requested id
func UpdateUser(uid string, user *models.User) (*models.User, error) {

	// To hash the password
	err := user.BeforeSave()
	if err != nil {
		return nil, err
	}
	DBClient = DBClient.Debug().Model(&models.User{}).Where("id = ?", uid).Take(&models.User{}).UpdateColumns(
		map[string]interface{}{
			"firstname":  user.Firstname,
			"lastname":   user.Lastname,
			"password":   user.Password,
			"passport":   user.Passport,
			"email":      user.Email,
			"updated_at": time.Now(),
		},
	)
	if DBClient.Error != nil {
		return nil, DBClient.Error
	}
	// This is the display the updated user
	err = DBClient.Debug().Model(&models.User{}).Where("id = ?", uid).Take(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
