package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

func main() {
	db, err = sql.Open("mysql", "root:@/golang-assessment")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
}
