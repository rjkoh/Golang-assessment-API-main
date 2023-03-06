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
}
