package database

import (
	"fmt"
	"forum/internal/config"
	"log"
)

func IsThereAnyCategories() bool {
	var	query	string
	var	c		config.CategoriesTableKeys
	var	count	int
	var	err		error

	c = config.TableKeys.Categories
	query = `
		SELECT COUNT(*) FROM `+c.Categories+`;
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
	var	c			config.CategoriesTableKeys
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

func AddPostCategories(postId int, categories_id ...int) {

}

func GetPostCategories(postId int) {

}
