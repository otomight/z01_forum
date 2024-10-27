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

	log.Println("Database initialized successfully")
	return nil
}

const schema = `
-- sessions table schema
CREATE TABLE IF NOT EXISTS Sessions (
    SessionID TEXT PRIMARY KEY,
    userID INTEGER,
    expiration TIMESTAMP,
    creationDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updateDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deletionDate TIMESTAMP,
    isDeleted BOOLEAN DEFAULT FALSE,
    FOREIGN KEY(UserID) REFERENCES clients(UserID)
);

-- posts table schema
CREATE TABLE IF NOT EXISTS Posts (
    postID INTEGER PRIMARY KEY AUTOINCREMENT,
    authorID INTEGER,
    Title TEXT,
    category TEXT,
    content TEXT,
    creationDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updateDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    DeletionDate TIMESTAMP,
    isDeleted BOOLEAN DEFAULT FALSE,
    FOREIGN KEY(AuthorID) REFERENCES clients(UserID) ON DELETE CASCADE
);

-- clients table schema
CREATE TABLE IF NOT EXISTS Clients (
    userID INTEGER PRIMARY KEY AUTOINCREMENT,
    lastName TEXT,
    firstName TEXT,
    userName TEXT UNIQUE,
    email TEXT UNIQUE,
    password TEXT,
    avatar TEXT,
    birthDate DATE,
    creationDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updateDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deletionDate TIMESTAMP 
);
`
