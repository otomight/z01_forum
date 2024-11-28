package database

import (
	"fmt"
	"forum/internal/config"
	"log"
)

func IsThereAnyCategories() bool {
	var	query	string
	var	k		config.StructTablesKeys = config.TableKeys
	var	count	int
	var	err		error

	query = `
		SELECT COUNT(*) FROM `+k.Categories.Table+`;
	`
	err = DB.QueryRow(query).Scan(&count)
	if err != nil {
		log.Printf("Error checking if there is any categorie in db: %v\n", err)
		return true
	}
	if count > 0 {
		return true
	}
	return false
}

// called on db init
func InsertCategories() {
	var	query		string
	var	categories	string
	var	err			error
	var	i			int

	if IsThereAnyCategories() {
		return
	}
	for i = 0; i < len(config.CategoriesNames); i++ {
		categories += fmt.Sprintf("('%s')", config.CategoriesNames[i])
		if i + 1 != len(config.CategoriesNames) {
			categories += " ,"
		}
	}
	query = `
		INSERT INTO categories (name)
		VALUES `+categories+`;
	`
	_, err = DB.Exec(query)
	if err != nil {
		log.Println("Categories not created:", err)
	}
}

func AddPostCategories(postId int, categories_id ...int) {

}

func GetPostCategories(postId int) {

}
