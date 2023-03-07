package routes

import (
	"github.com/gorilla/mux"
	"github.com/rjkoh/golang-assessment-api/pkg/controllers"
)

func APIRoutes(router *mux.Router) {
	router.HandleFunc("/api/register", controllers.RegisterStudents).Methods("POST")
	router.HandleFunc("/api/commonstudents", controllers.FindCommonStudents).
		Queries("teacher", "{teacher:.*}").
		Methods("GET")
	router.HandleFunc("/api/suspend", controllers.SuspendStudent).Methods("POST")
	router.HandleFunc("/api/retrievefornotifications", controllers.RetrieveStudentsForNotification).Methods("POST")
}
