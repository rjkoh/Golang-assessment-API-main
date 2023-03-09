package config

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

const password = ""

func Connect() {
	db, err = sql.Open("mysql", fmt.Sprintf("root:%s@tcp(localhost:3306)/", password))
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS testdb")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("USE testdb")
	if err != nil {
		panic(err)
	}
}

func GetDB() *sql.DB {
	return db
}
