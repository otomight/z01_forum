package database

import (
	"database/sql"
	"fmt"
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

func isThereAnyCategories() bool {
	var	c		config.CategoriesTableKeys

	c = config.TableKeys.Categories
	if countRows(c.Categories) > 0 {
		return true
	}
	return false
}

// called on db init
func insertCategories() {
	var	query		string
	var	categories	string
	var	err			error
	var	c			config.CategoriesTableKeys
	var	i			int

	if isThereAnyCategories() {
		return
	}
	for i = 0; i < len(config.CategoriesNames); i++ {
		categories += fmt.Sprintf("('%s')", config.CategoriesNames[i])
		if i + 1 != len(config.CategoriesNames) {
			categories += " ,"
		}
	}
	c = config.TableKeys.Categories
	query = `
		INSERT INTO `+c.Categories+` (`+c.Name+`)
		VALUES `+categories+`;
	`
	_, err = DB.Exec(query)
	if err != nil {
		log.Println("Categories not created:", err)
	}
}

func InitDB() error {
	//Open DB connection
	var	schema	string
	var	err		error

	DB, err = sql.Open("sqlite3", config.DbFilePath)
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

	insertCategories()
	InsertSampleClient()
	InsertSamplePost()

	log.Println("Database initialized successfully")
	return nil
}
