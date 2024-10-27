package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() error {
	//Open DB connection
	var err error
	DB, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}

	//Enable foreign key for SQLite
	if _, err := DB.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return err
	}

	//Execute create_table schema script
	if _, err := DB.Exec(schema); err != nil {
		return err
	}

	InsertSampleClient()
	InsertSamplePost()

	log.Println("Database initialized successfully")
	return nil
}

const schema = `
-- sessions table schema
CREATE TABLE IF NOT EXISTS Sessions (
    session_id TEXT PRIMARY KEY,
    user_id INTEGER,
    expiration TIMESTAMP,
    creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deletion_date TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE,
    FOREIGN KEY(user_id) REFERENCES clients(user_id)
);

-- posts table schema
CREATE TABLE IF NOT EXISTS Posts (
    post_id INTEGER PRIMARY KEY AUTOINCREMENT,
    author_id INTEGER,
    title TEXT,
    category TEXT,
    content TEXT,
    creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deletion_date TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE,
    FOREIGN KEY(author_id) REFERENCES clients(user_id) ON DELETE CASCADE
);

-- clients table schema
CREATE TABLE IF NOT EXISTS Clients (
    user_id INTEGER PRIMARY KEY AUTOINCREMENT,
    last_name TEXT,
    first_name TEXT,
    user_name TEXT UNIQUE,
    email TEXT UNIQUE,
    password TEXT,
    avatar TEXT,
    birth_date DATE,
    creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deletion_date TIMESTAMP 
);
`
