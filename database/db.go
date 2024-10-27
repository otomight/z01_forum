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

	//Execute create_table schema script
	if _, err := DB.Exec(schema); err != nil {
		log.Fatal("Error executing schema:", err)
	}
	log.Println("Database initialized successfully")
	return DB, nil
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
    FOREIGN KEY(user_id) REFERENCES clients(user_id)
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
    FOREIGN KEY(author_id) REFERENCES clients(user_id)
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
