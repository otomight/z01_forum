package database

import (
	"database/sql"
	"fmt"
	"forum/internal/config"
	"log"
)

func AddPostCategories(postID int, categoriesID ...int) error {
	var	query	string
	var	pc		config.PostsCategoriesTableKeys
	var	args	[]any
	var	i		int
	var	err		error

	if len(categoriesID) == 0 {
		return fmt.Errorf("No categories id provided")
	}
	pc = config.TableKeys.PostsCategories
	query = `
		INSERT INTO `+pc.PostsCategories+` (
			`+pc.PostID+`, `+pc.CategoryID+`
		)
		VALUES
	`
	args = make([]any, len(categoriesID)*2)
	for i = 0; i < len(categoriesID); i++ {
		query += " (?, ?)"
		args[i*2] = postID
		args[i*2 + 1] = categoriesID[i]
		if i + 1 != len(categoriesID) {
			query += ","
		} else {
			query += ";"
		}
	}
	_, err = DB.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func GetPostCategories(postID int) ([]*Category, error) {
	var	query			string
	var	pc				config.PostsCategoriesTableKeys
	var	c				config.CategoriesTableKeys
	var	rows			*sql.Rows
	var	postCategories	[]*Category
	var	postCategory	*Category
	var	err				error

	c = config.TableKeys.Categories
	pc = config.TableKeys.PostsCategories
	query = `
		SELECT c.`+c.ID+`, c.`+c.Name+`
		FROM `+c.Categories+` c
		JOIN `+pc.PostsCategories+` pc ON pc.`+pc.CategoryID+` = c.`+c.ID+`
		WHERE pc.`+pc.PostID+` = ?;
	`
	rows, err = DB.Query(query, postID)
	if err != nil {
		return []*Category{}, err
	}
	defer rows.Close()
	for rows.Next() {
		postCategory = &Category{}
		err = rows.Scan(&postCategory.ID, &postCategory.Name)
		if err != nil {
			log.Printf("Unexpected error at scan category name: %v", err)
		}
		postCategories = append(postCategories, postCategory)
	}
	return postCategories, nil
}

func GetGlobalCategories() ([]*Category, error) {
	var	query		string
	var	c			config.CategoriesTableKeys
	var	categories	[]*Category
	var	category	*Category
	var	rows		*sql.Rows
	var	err			error

	c = config.TableKeys.Categories
	query = `
		SELECT `+c.ID+`, `+c.Name+`
		FROM `+c.Categories+`
	`
	rows, err = DB.Query(query)
	if err != nil {
		log.Printf("Error fetching categories from database: %v", err)
		return []*Category{}, err
	}
	defer rows.Close()
	for rows.Next() {
		category = &Category{}
		err = rows.Scan(&category.ID, &category.Name)
		if err != nil {
			log.Printf("Unexpected error at scan category name: %v", err)
		}
		categories = append(categories, category)
	}
	return categories, nil
}
