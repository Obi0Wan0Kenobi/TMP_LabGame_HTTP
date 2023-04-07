package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var DB *sql.DB

func ConnectToDB() {
	var err error
	DB, err = sql.Open("sqlite3", "db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
}
