package database

import "forum/internal/config"

func AddPostCategories(postID int, categoriesID ...int) error {
	var	query	string
	var	pc		config.PostsCategoriesTableKeys
	var	args	[]any
	var	i		int
	var	err		error

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
	_, err = DB.Exec(query, postID, categoriesID)
	if err != nil {
		return err
	}
	return nil
}

func GetPostCategories(postID int) {

}
