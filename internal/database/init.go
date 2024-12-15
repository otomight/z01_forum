package database

import (
	"database/sql"
	"forum/internal/config"
	"forum/internal/utils"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func getSqlTables() string {
	var content []byte
	var err error

	content, err = os.ReadFile(config.SqlTablesFilePath)
	if err != nil {
		log.Println(".sql file not found, .db has not been created.")
	}
	return string(content)
}

func isThereAnyCategories() bool {
	var c config.CategoriesTableKeys

	c = config.TableKeys.Categories
	if countRows(c.Categories) > 0 {
		return true
	}
	return false
}

// called on db init
func insertCategories() {
	var err error
	var c config.CategoriesTableKeys

	if isThereAnyCategories() {
		return
	}
	c = config.TableKeys.Categories
	_, err = insertInto(InsertIntoQuery{
		Table:  c.Categories,
		Keys:   []string{c.Name},
		Values: utils.ToMatrix(config.CategoriesNames),
	})
	if err != nil {
		log.Println("Categories not created:", err)
	}
}

func InitDB() error {
	//Open DB connection
	var schema string
	var err error

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
