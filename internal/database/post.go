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
		Keys:	[]string{p.AuthorID, p.Title, p.Content, p.ImagePath},
		Values:	[][]any{{post.AuthorID, post.Title, post.Content, post.ImagePath}},
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
	curUserID int, condition string, args ...any,
) (*sql.Rows, error) {
	var	p		config.PostsTableKeys
	var	cl		config.ClientsTableKeys
	var	pc		config.PostsCategoriesTableKeys
	var	pr		config.PostsReactionsTableKeys
	var	c		config.CommentsTableKeys
	var	query	string
	var	rows	*sql.Rows
	var err		error

	p = config.TableKeys.Posts
	cl = config.TableKeys.Clients
	pc = config.TableKeys.PostsCategories
	pr = config.TableKeys.PostsReactions
	c = config.TableKeys.Comments
	query = `
		SELECT DISTINCT p.`+p.ID+`, p.`+p.AuthorID+`, cl.`+cl.UserName+`,
			p.`+p.Title+`, p.`+p.Content+`, p.`+p.ImagePath+`,
			p.`+p.CreationDate+`, p.`+p.UpdateDate+`,
			p.`+p.Likes+`, p.`+p.Dislikes+`, upr.`+pr.Liked+`
		FROM `+p.Posts+` p
		JOIN `+cl.Clients+` cl ON p.`+p.AuthorID+` = cl.`+cl.ID+`
		LEFT JOIN `+pr.PostsReactions+` upr
			ON upr.`+pr.PostID+` = p.`+p.ID+` AND upr.`+pr.UserID+` = ?
		LEFT JOIN `+pr.PostsReactions+` pr ON pr.`+pr.PostID+` = p.`+p.ID+`
		LEFT JOIN `+pc.PostsCategories+` pc ON pc.`+pc.PostID+` = p.`+p.ID+`
		LEFT JOIN `+c.Comments+` c ON c.`+c.PostID+` = p.`+p.ID+`
	` // upr is for user post reaction, it fetch the curUserID reaction
	if condition != "" {
		query += ` WHERE `+condition+``
	}
	query += ";"
	rows, err = DB.Query(query, append([]any{curUserID}, args...)...)
	return rows, err
}

func fillPostExternalData(
	curUserID int, post *Post, userLiked *bool, includeComments bool,
) {
	var	err	error

	post.UserConfig = getUserConfig(userLiked)
	if includeComments {
		post.Comments, err = GetCommentsByPostID(curUserID, post.ID)
		if err != nil {
			log.Println(err.Error())
		}
	}
	post.Categories, err = GetPostCategories(post.ID)
	if err != nil {
		log.Println(err.Error())
	}
}

func getPostsWithCondition(
	curUserID int, includeComments bool, condition string, args ...any,
) ([]*Post, error) {
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
			&post.Content, &post.ImagePath, &post.CreationDate,
			&post.UpdateDate, &post.Likes, &post.Dislikes, &userLiked,
		)
		if err != nil {
			log.Printf("Error scanning post: %v\n", err)
			continue
		}
		fillPostExternalData(curUserID, post, userLiked, includeComments)
		posts = append(posts, post)
	}
	if err = rows.Err(); err != nil {
		log.Println("Error during row iteration")
		return nil, err
	}
	return posts, nil
}

func GetPostByID(curUserID int, postID int) (*Post, error) {
	var	condition	string
	var	p			config.PostsTableKeys
	var	posts		[]*Post
	var	err			error

	p = config.TableKeys.Posts
	condition = `p.`+p.ID+` = ?`
	posts, err = getPostsWithCondition(curUserID, true, condition, postID)
	if len(posts) == 0 {
		return nil, err
	}
	return posts[0], err
}

func GetAllPosts(curUserID int) ([]*Post, error) {
	return getPostsWithCondition(curUserID, true, "")
}

// For minimal post
func GetPostsCreatedByUser(curUserID int, userID int) ([]*Post, error) {
	var	condition	string
	var	p			config.PostsTableKeys

	p = config.TableKeys.Posts
	condition = `p.`+p.AuthorID+` = ?`
	return getPostsWithCondition(curUserID, false, condition, userID)
}

// For minimal post
func GetPostsLikedByUser(curUserID int, userID int) ([]*Post, error) {
	var	condition	string
	var	pr			config.PostsReactionsTableKeys

	pr = config.TableKeys.PostsReactions
	condition = `pr.`+pr.UserID+` = ? AND pr.`+pr.Liked+` = true`
	return getPostsWithCondition(curUserID, false, condition, userID)
}

// For minimal post
func GetPostsDislikedByUser(curUserID int, userID int) ([]*Post, error) {
	var	condition	string
	var	pr			config.PostsReactionsTableKeys

	pr = config.TableKeys.PostsReactions
	condition = `pr.`+pr.UserID+` = ? AND pr.`+pr.Liked+` = false`
	return getPostsWithCondition(curUserID, false, condition, userID)
}

// For minimal post
func GetPostsCommentedByUser(curUserID int, userID int) ([]*Post, error) {
	var	condition	string
	var	posts		[]*Post
	var	c			config.CommentsTableKeys
	var	i			int
	var	err			error

	c = config.TableKeys.Comments
	condition = `c.`+c.UserID+` = ?`
	posts, err = getPostsWithCondition(curUserID, false, condition, userID)
	if err != nil {
		return posts, err
	}
	for i = 0; i < len(posts); i++ {
		posts[i].Comments, err = GetCommentsOfPostFromUser(
			curUserID, posts[i].ID, curUserID,
		)
		if err != nil {
			log.Printf("Error at fetching comments of post %d\n", posts[i].ID)
		}
	}
	return posts, err
}

// For minimal post
func GetPostsByCategoryID(curUserID int, categoryID int) ([]*Post, error) {
	var	condition	string
	var	pc			config.PostsCategoriesTableKeys

	pc = config.TableKeys.PostsCategories
	condition = `pc.`+pc.CategoryID+` = ?`
	return getPostsWithCondition(curUserID, false, condition, categoryID)
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
