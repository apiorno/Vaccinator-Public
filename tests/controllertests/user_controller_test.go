package controllertests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/S-S-Group/Vaccinator/src/controllers"
	"github.com/S-S-Group/Vaccinator/src/models"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/assert.v1"
)

func TestCreateUser(t *testing.T) {
	userController := &controllers.UsersController{}
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	samples := []struct {
		inputJSON    string
		statusCode   int
		firstname    string
		lastname     string
		passport     string
		email        string
		errorMessage string
	}{
		{
			inputJSON:    `{"firstname":"Pet", "lastname":"da","passport": "1","email": "pet@gmail.com", "password": "password"}`,
			statusCode:   201,
			firstname:    "Pet",
			lastname:     "da",
			passport:     "1",
			email:        "pet@gmail.com",
			errorMessage: "",
		},
		{
			inputJSON:    `{"firstname":"Frank", "lastname":"de","passport": "2","email": "pet@gmail.com", "password": "password"}`,
			statusCode:   500,
			errorMessage: "Email Already Taken",
		},
		{
			inputJSON:    `{"firstname":"Pet", "lastname":"di","passport": "3","email": "grand@gmail.com", "password": "password"}`,
			statusCode:   500,
			errorMessage: "Firstname Already Taken",
		},
		{
			inputJSON:    `{"firstname":"Fi", "lastname":"da","passport": "4","email": "grand@gmail.com", "password": "password"}`,
			statusCode:   500,
			errorMessage: "Lastname Already Taken",
		},
		{
			inputJSON:    `{"firstname":"Frank", "lastname":"de","passport": "1","email": "cat@gmail.com", "password": "password"}`,
			statusCode:   500,
			errorMessage: "Passport Already Taken",
		},
		{
			inputJSON:    `{"firstname":"Kan", "lastname":"do","passport": "5", "email": "kangmail.com", "password": "password"}`,
			statusCode:   422,
			errorMessage: "Invalid Email",
		},
		{
			inputJSON:    `{"firstname": "", "lastname":"du","passport": "6","email": "kan@gmail.com", "password": "password"}`,
			statusCode:   422,
			errorMessage: "Required Firstname",
		},
		{
			inputJSON:    `{"firstname": "firstname", "lastname":"","passport": "7","email": "kan@gmail.com", "password": "password"}`,
			statusCode:   422,
			errorMessage: "Required Lastname",
		},
		{
			inputJSON:    `{"firstname": "Kan", "lastname":"er","passport": "8","email": "", "password": "password"}`,
			statusCode:   422,
			errorMessage: "Required Email",
		},
		{
			inputJSON:    `{"firstname": "Kan", "lastname":"or","passport": "9","email": "kan@gmail.com", "password": ""}`,
			statusCode:   422,
			errorMessage: "Required Password",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/users", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(userController.CreateUser())
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["firstname"], v.firstname)
			assert.Equal(t, responseMap["lastname"], v.lastname)
			assert.Equal(t, responseMap["passport"], v.passport)
			assert.Equal(t, responseMap["email"], v.email)
		}
		if v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetUsers(t *testing.T) {
	userController := &controllers.UsersController{}
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	_, err = seedUsers()
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(userController.GetAllUsers())
	handler.ServeHTTP(rr, req)

	var users []models.User
	err = json.Unmarshal([]byte(rr.Body.String()), &users)
	if err != nil {
		log.Fatalf("Cannot convert to json: %v\n", err)
	}
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(users), 2)
}

func TestGetUserByID(t *testing.T) {
	userController := &controllers.UsersController{}
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	user, err := seedOneUser()
	if err != nil {
		log.Fatal(err)
	}
	userSample := []struct {
		id           string
		statusCode   int
		firstname    string
		lastname     string
		passport     string
		email        string
		errorMessage string
	}{
		{
			id:         strconv.Itoa(int(user.ID)),
			statusCode: 200,
			firstname:  user.Firstname,
			lastname:   user.Lastname,
			passport:   user.Passport,
			email:      user.Email,
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
	}
	for _, v := range userSample {

		req, err := http.NewRequest("GET", "/users", nil)
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(userController.GetUserById())
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, user.Firstname, responseMap["firstname"])
			assert.Equal(t, user.Lastname, responseMap["lastname"])
			assert.Equal(t, user.Passport, responseMap["passport"])
			assert.Equal(t, user.Email, responseMap["email"])
		}
	}
}

func TestUpdateUser(t *testing.T) {

	var AuthEmail, AuthPassword string
	var AuthID uint32

	userController := &controllers.UsersController{}
	loginController := &controllers.LoginController{}
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	users, err := seedUsers() //we need atleast two users to properly check the update
	if err != nil {
		log.Fatalf("Error seeding user: %v\n", err)
	}
	// Get only the first user
	for _, user := range users {
		if user.ID == 2 {
			continue
		}
		AuthID = user.ID
		AuthEmail = user.Email
		AuthPassword = "password" //Note the password in the database is already hashed, we want unhashed
	}
	//Login the user and get the authentication token
	token, err := loginController.SignIn(AuthEmail, AuthPassword, 1)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		id              string
		updateJSON      string
		statusCode      int
		updateFirstname string
		updateLastname  string
		updatePassport  string
		updateEmail     string
		tokenGiven      string
		errorMessage    string
	}{
		{
			// Convert int32 to int first before converting to string
			id:              strconv.Itoa(int(AuthID)),
			updateJSON:      `{"firstname":"Grand", "lastname": "da", "passport": "1","email": "grand@gmail.com", "password": "password"}`,
			statusCode:      200,
			updateFirstname: "Grand",
			updateLastname:  "da",
			updatePassport:  "1",
			updateEmail:     "grand@gmail.com",
			tokenGiven:      tokenString,
			errorMessage:    "",
		},
		{
			// When password field is empty
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"firstname":"Woman", "lastname": "de", "passport": "2", "email": "woman@gmail.com", "password": ""}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Password",
		},
		{
			// When no token was passed
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"firstname":"Man", "lastname": "di", "passport": "3","email": "man@gmail.com", "password": "password"}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token was passed
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"firstname":"Woman", "lastname": "do", "passport": "4","email": "woman@gmail.com", "password": "password"}`,
			statusCode:   401,
			tokenGiven:   "This is incorrect token",
			errorMessage: "Unauthorized",
		},
		{
			// Remember "kenny@gmail.com" belongs to user 2
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"firstname":"Frank", "lastname": "du", "passport": "5","email": "kenny@gmail.com", "password": "password"}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "Email Already Taken",
		},
		{
			// Remember "Kenny Morris" belongs to user 2
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"firstname":"Kenny Morris", "lastname": "er", "passport": "6","email": "grand@gmail.com", "password": "password"}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "Firstname Already Taken",
		},
		{
			// Remember "er" belongs to user 6
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"firstname":"Kenny", "lastname": "er", "passport": "10","email": "kiki@gmail.com", "password": "password"}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "Lastname Already Taken",
		},
		{
			// Remember "er" belongs to user 6
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"firstname":"Kanny", "lastname": "ru", "passport": "6","email": "kiri@gmail.com", "password": "password"}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "Passport Already Taken",
		},
		{
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"firstname":"Kan", "lastname": "ir", "passport": "7","email": "kangmail.com", "password": "password"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Invalid Email",
		},
		{
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"firstname": "", "lastname": "or", "passport": "8","email": "kan@gmail.com", "password": "password"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Firstname",
		},
		{
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"firstname": "asd", "lastname": "", "passport": "11","email": "kin@gmail.com", "password": "password"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Lastname",
		},
		{
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"firstname": "sda", "lastname": "rar", "passport": "","email": "kon@gmail.com", "password": "password"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Passport",
		},
		{
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"firstname": "Kan", "lastname": "ar", "passport": "9","email": "", "password": "password"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Email",
		},
		{
			id:         "unknwon",
			tokenGiven: tokenString,
			statusCode: 400,
		},
		{
			// When user 2 is using user 1 token
			id:           strconv.Itoa(int(2)),
			updateJSON:   `{"firstname": "Mike","lastname": "nike","passport": "12" "email": "mike@gmail.com", "password": "password"}`,
			tokenGiven:   tokenString,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/users", bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(userController.UpdateUser())

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, responseMap["firstname"], v.updateFirstname)
			assert.Equal(t, responseMap["lastname"], v.updateLastname)
			assert.Equal(t, responseMap["passport"], v.updatePassport)
			assert.Equal(t, responseMap["email"], v.updateEmail)
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestDeleteUser(t *testing.T) {

	var AuthEmail, AuthPassword string
	var AuthID uint32
	userController := &controllers.UsersController{}
	loginController := &controllers.LoginController{}
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	users, err := seedUsers() //we need atleast two users to properly check the update
	if err != nil {
		log.Fatalf("Error seeding user: %v\n", err)
	}
	// Get only the first and log him in
	for _, user := range users {
		if user.ID == 2 {
			continue
		}
		AuthID = user.ID
		AuthEmail = user.Email
		AuthPassword = "password" ////Note the password in the database is already hashed, we want unhashed
	}
	//Login the user and get the authentication token
	token, err := loginController.SignIn(AuthEmail, AuthPassword, 1)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	userSample := []struct {
		id           string
		tokenGiven   string
		statusCode   int
		errorMessage string
	}{
		{
			// Convert int32 to int first before converting to string
			id:           strconv.Itoa(int(AuthID)),
			tokenGiven:   tokenString,
			statusCode:   204,
			errorMessage: "",
		},
		{
			// When no token is given
			id:           strconv.Itoa(int(AuthID)),
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is given
			id:           strconv.Itoa(int(AuthID)),
			tokenGiven:   "This is an incorrect token",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			id:         "unknwon",
			tokenGiven: tokenString,
			statusCode: 400,
		},
		{
			// User 2 trying to use User 1 token
			id:           strconv.Itoa(int(2)),
			tokenGiven:   tokenString,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range userSample {

		req, err := http.NewRequest("GET", "/users", nil)
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(userController.DeleteUser())

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 401 && v.errorMessage != "" {
			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
