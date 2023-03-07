package models

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rjkoh/golang-assessment-api/pkg/config"
)

var db *sql.DB

func init() {
	config.Connect()
	db = config.GetDB()
	defer db.Close()

	_, err := db.Exec("CREATE TABLE IF NOT EXISTS student (Email TEXT NOT NULL, Teacher TEXT NOT NULL, IsSuspended BOOLEAN NOT NULL DEFAULT FALSE)")
	if err != nil {
		panic(err)
	}
}
