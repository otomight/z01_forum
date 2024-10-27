package database

import (
	"database/sql"
	"log"
)

var DB *sql.DB

func InitDB() (*sql.DB, error) {
	//Open DB connection
	var err error
	DB, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}

	//Enable foreign key for SQLite
	_, err = DB.Exec("PRAGMA foreign_key = ON")
	if err != nil {
		return nil, err
	}

	log.Println("Database initialized successfully")
	return DB, nil
}
