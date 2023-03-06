package routes

import (
	"github.com/gorilla/mux"
	"github.com/rjkoh/golang-assessment-api/controllers"
)

func APIRoutes(router *mux.Router) {
	router.HandleFunc("/api/register", controllers.RegisterStudents).Methods("POST")
	router.HandleFunc("/api/commonstudents/{teacher+}", controllers.findCommonStudents).Methods("GET")
	router.HandleFunc("/api/suspend", controllers.SuspendStudent).Methods("POST")
	router.HandleFunc("/api/retrievefornotifications", controllers.SuspendStudent).Methods("POST")
}
