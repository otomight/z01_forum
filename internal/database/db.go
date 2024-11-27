package database

import (
	"database/sql"
	"forum/internal/config"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func getSqlTables() string {
	var	content	[]byte
	var	err		error

	content, err = os.ReadFile(config.SqlTablesFilePath)
	if err != nil {
		log.Println(".sql file not found, .db has not been created.")
	}
	return string(content)
}

func InitDB() error {
	//Open DB connection
	var	schema	string
	var	err		error

	DB, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}

	//Enable foreign key for SQLite
	if _, err := DB.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return err
	}

	//Execute create_table schema script
	schema = getSqlTables()
	if _, err := DB.Exec(schema); err != nil {
		return err
	}

	InsertSampleClient()
	InsertSamplePost()

	log.Println("Database initialized successfully")
	return nil
}
