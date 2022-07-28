package controllers

import (
	"log"
	"net/http"

	"github.com/S-S-Group/Vaccinator/src/middlewares"
	"github.com/gorilla/mux"
)

func InitializeUsersController(mux *mux.Router, l *log.Logger) {
	usersController := &UsersController{l}

	getRouter := mux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(usersController.GetUserById()))
	getRouter.HandleFunc("/users", middlewares.SetMiddlewareJSON(usersController.GetAllUsers()))

	postRouter := mux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/users", middlewares.SetMiddlewareJSON(usersController.CreateUser()))

	deleteRouter := mux.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(usersController.DeleteUser()))

	putRouter := mux.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(usersController.UpdateUser())))
}

func InitializeNotificationsController(mux *mux.Router, l *log.Logger) {
	notificationsController := &NotificationsController{l}

	getRouter := mux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/notifications/{id}", middlewares.SetMiddlewareJSON(notificationsController.GetNotificationById()))
	getRouter.HandleFunc("/notifications", middlewares.SetMiddlewareJSON(notificationsController.GetAllNotifications()))

	postRouter := mux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/notifications", middlewares.SetMiddlewareJSON(notificationsController.CreateNotification()))

	deleteRouter := mux.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/notifications/{id}", middlewares.SetMiddlewareAuthentication(notificationsController.DeleteNotification()))

	putRouter := mux.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/notifications/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(notificationsController.UpdateNotification())))
}

func InitializeLoginController(mux *mux.Router, l *log.Logger) {
	loginController := &LoginController{l}
	postRouter := mux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/login", middlewares.SetMiddlewareJSON(loginController.Login))
}
func InitializeCertificationsController(mux *mux.Router, l *log.Logger) {
	certificationController := &CertificationsController{l}

	postRouter := mux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/certifications", certificationController.CreateCertification())

	getRouter := mux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/certifications/{id}", certificationController.GetCertificationsOfUser())

	updateRouter := mux.Methods(http.MethodPut).Subrouter()
	updateRouter.HandleFunc("/certifications/{id}", certificationController.UpdateCertification())

	deleteRouter := mux.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/certifications/{id}", certificationController.DeleteCertification())
}
func InitializeAssistancesController(mux *mux.Router, l *log.Logger) {
	assistancesController := &AssistancesController{l}

	countRouter := mux.Methods(http.MethodGet).Subrouter()
	countRouter.HandleFunc("/assistances/count", assistancesController.GetAssistancesCount())

	postRouter := mux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/assistances", assistancesController.CreateAssistance())

	getRouter := mux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/assistances/{id}", assistancesController.GetAssistancesOfUser())

	updateRouter := mux.Methods(http.MethodPut).Subrouter()
	updateRouter.HandleFunc("/assistances/{id}", assistancesController.UpdateAssistance())

	deleteRouter := mux.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/assistances/{id}", assistancesController.DeleteAssistance())
}
func InitializeValidationsController(mux *mux.Router, l *log.Logger) {
	validationController := &ValidationsController{l}

	postRouter := mux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/validations", validationController.ValidateUserCertifications())
}

// Startup initializes the controller
func Startup(l *log.Logger) *mux.Router {
	newMux := mux.NewRouter()
	InitializeUsersController(newMux, l)
	InitializeNotificationsController(newMux, l)
	InitializeCertificationsController(newMux, l)
	InitializeLoginController(newMux, l)
	InitializeValidationsController(newMux, l)
	InitializeAssistancesController(newMux, l)

	return newMux
}
