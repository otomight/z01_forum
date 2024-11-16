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
    user_role TEXT DEFAULT 'user',
    creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deletion_date TIMESTAMP 
);

-- sessions table schema
CREATE TABLE IF NOT EXISTS Sessions (
    session_id TEXT PRIMARY KEY,
    user_id INTEGER,
    user_role TEXT,
    is_logged_in BOOLEAN DEFAULT FALSE,
    expiration TIMESTAMP,
    creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deletion_date TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE,
    FOREIGN KEY(user_id) REFERENCES clients(user_id),
    FOREIGN kEY (user_name) REFERENCES Client(user_name)
);

-- posts table schema
CREATE TABLE IF NOT EXISTS Posts (
    post_id INTEGER PRIMARY KEY AUTOINCREMENT,
    author_id INTEGER,
    title TEXT,
    category TEXT,
    tags TEXT,
    content TEXT,
    creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deletion_date TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE,
    FOREIGN KEY(author_id) REFERENCES clients(user_id) ON DELETE CASCADE
);

-- comments table schema
CREATE TABLE IF NOT EXISTS Comments (
    comment_id INT AUTO_INCREMENT PRIMARY KEY, 
    post_id INT NOT NULL,
    user_id INT NOT NULL,
    content TEXT,
    creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES Posts(post_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES Clients(user_id) ON DELETE CASCADE
);
`
