package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rjkoh/golang-assessment-api/pkg/routes"
)

func main() {
	router := mux.NewRouter()
	routes.APIRoutes(router)
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":9010", router))
}
