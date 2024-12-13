package database

import (
	"database/sql"
	"fmt"
	"forum/internal/config"
	"log"
)

// Post CRUD operations
func NewPost(post *Post, categoriesIDs []int) (int, error) {
	var	p		config.PostsTableKeys

	if len(categoriesIDs) == 0 {
		return 0, fmt.Errorf("No categories provided.")
	}
	p = config.TableKeys.Posts
	result, err := insertInto(InsertIntoQuery{
		Table:	p.Posts,
		Keys:	[]string{p.AuthorID, p.Title, p.Content},
		Values:	[][]any{{post.AuthorID, post.Title, post.Content}},
	})
	if err != nil {
		return 0, err
	}
	postID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	if err = AddPostCategories(int(postID), categoriesIDs...); err != nil {
		log.Printf("Error adding categories to the database: %v", err)
	}
	return int(postID), err
}

func getPostsWithConditionQueryResult(
	userID int, condition string, args ...any,
) (*sql.Rows, error) {
	var	p			config.PostsTableKeys
	var	cl			config.ClientsTableKeys
	var	pc			config.PostsCategoriesTableKeys
	var	pr			config.PostsReactionsTableKeys
	var	query		string
	var	rows		*sql.Rows
	var err			error

	p = config.TableKeys.Posts
	cl = config.TableKeys.Clients
	pc = config.TableKeys.PostsCategories
	pr = config.TableKeys.PostsReactions
	query = `
		SELECT DISTINCT p.`+p.ID+`, p.`+p.AuthorID+`, cl.`+cl.UserName+`,
			p.`+p.Title+`, p.`+p.Content+`, p.`+p.CreationDate+`,
			p.`+p.UpdateDate+`, p.`+p.Likes+`, p.`+p.Dislikes+`, pr.`+pr.Liked+`
		FROM `+p.Posts+` p
		JOIN `+cl.Clients+` cl ON p.`+p.AuthorID+` = cl.`+cl.ID+`
		LEFT JOIN `+pc.PostsCategories+` pc ON pc.`+pc.PostID+` = p.`+p.ID+`
		LEFT JOIN `+pr.PostsReactions+` pr
		ON pr.`+pr.PostID+` = p.`+p.ID+` AND pr.`+pr.UserID+` = ?
	`
	if condition != "" {
		query += ` WHERE `+condition+``
	}
	query += ";"
	rows, err = DB.Query(query, append([]any{userID}, args...)...)
	return rows, err
}

func fillPostExternalData(curUserID int, post *Post, userLiked *bool) {
	var	err	error

	post.UserConfig = getUserConfig(userLiked)
	post.Comments, err = GetCommentsByPostID(curUserID, post.ID)
	if err != nil {
		log.Println(err.Error())
	}
	post.Categories, err = GetPostCategories(post.ID)
	if err != nil {
		log.Println(err.Error())
	}
}

func getPostsWithCondition(curUserID int, condition string, args ...any) ([]*Post, error) {
	var	posts		[]*Post
	var	post		*Post
	var	rows		*sql.Rows
	var	userLiked	*bool
	var	err			error

	rows, err = getPostsWithConditionQueryResult(curUserID, condition, args...)
	if err != nil {
		log.Println("Error on post query")
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		post = &Post{}
		err = rows.Scan(
			&post.ID, &post.AuthorID, &post.UserName, &post.Title,
			&post.Content, &post.CreationDate, &post.UpdateDate,
			&post.Likes, &post.Dislikes, &userLiked,
		)
		if err != nil {
			log.Println("Error scanning post")
			continue
		}
		fillPostExternalData(curUserID, post, userLiked)
		posts = append(posts, post)
	}
	if err = rows.Err(); err != nil {
		log.Println("Error during row iteration")
		return nil, err
	}
	return posts, nil
}

func GetPostByID(curUserID int, userID int) (*Post, error) {
	var	condition	string
	var	p			config.PostsTableKeys
	var	posts		[]*Post
	var	err			error

	p = config.TableKeys.Posts
	condition = `p.`+p.ID+` = ?`
	posts, err = getPostsWithCondition(curUserID, condition, userID)
	if len(posts) == 0 {
		return nil, err
	}
	return posts[0], err
}

func GetPostsCreatedByUser(curUserID int, userID int) ([]*Post, error) {
	var	condition	string
	var	p			config.PostsTableKeys

	p = config.TableKeys.Posts
	condition = `p.`+p.AuthorID+` = ?`
	return getPostsWithCondition(curUserID, condition, userID)
}

func GetPostsLikedByUser(curUserID int, userID int) ([]*Post, error) {
	var	condition	string
	var	pr			config.PostsReactionsTableKeys

	pr = config.TableKeys.PostsReactions
	condition = ``+pr.UserID+` = ?`
	return getPostsWithCondition(curUserID, condition, userID)
}

func GetPostsRelatedToCurUser(userID int) ([]*Post, error) {
	var	condition	string
	var	p			config.PostsTableKeys
	var	r			config.PostsReactionsTableKeys

	p = config.TableKeys.Posts
	r = config.TableKeys.PostsReactions
	condition = `p.`+p.AuthorID+` = ? OR `+r.UserID+` = ?`
	return getPostsWithCondition(userID, condition, userID, userID)
}

func GetPostsByCategoryID(curUserID int, categoryID int) ([]*Post, error) {
	var	condition	string
	var	pc			config.PostsCategoriesTableKeys

	pc = config.TableKeys.PostsCategories
	condition = `pc.`+pc.CategoryID+` = ?`
	return getPostsWithCondition(curUserID, condition, categoryID)
}

// Retrieve all the posts from database
func GetAllPosts(curUserID int) ([]*Post, error) {
	return getPostsWithCondition(curUserID, "")
}

func deletePostWithCondition(condition string, args ...any) error {
	var	query	string
	var	p		config.PostsTableKeys
	var	err		error

	p = config.TableKeys.Posts
	query = `
		DELETE FROM `+p.Posts+`
	`
	if condition != "" {
		query += ` WHERE `+condition+``
	}
	query += ";"
	_, err = DB.Exec(query, args...)
	if err != nil {
		log.Println("Error deleting post")
		return err
	}
	return nil
}

func DeletePost(postID int) error {
	var	p			config.PostsTableKeys
	var	condition	string

	p = config.TableKeys.Posts
	condition = ``+p.ID+` = ?`
	return deletePostWithCondition(condition, postID)
}
