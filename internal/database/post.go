package database

import (
	"database/sql"
	"fmt"
	"forum/internal/config"
	"log"
)

// Post CRUD operations
func NewPost(post *Post, categoriesIDs []int) (int, error) {
	var	query	string
	var	p		config.PostsTableKeys

	p = config.TableKeys.Posts
	query = `
		INSERT INTO `+p.Posts+` (
			`+p.AuthorID+`, `+p.Title+`, `+p.Content+`,
			`+p.CreationDate+`, `+p.UpdateDate+`, `+p.IsDeleted+`
		)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 0)
	`
	result, err := DB.Exec(query, post.AuthorID, post.Title, post.Content)
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
	condition string, args ...any,
) (*sql.Rows, error) {
	var	p			config.PostsTableKeys
	var	cl			config.ClientsTableKeys
	var	pc			config.PostsCategoriesTableKeys
	var	query		string
	var	rows		*sql.Rows
	var	conditions	string
	var err			error

	p = config.TableKeys.Posts
	cl = config.TableKeys.Clients
	pc = config.TableKeys.PostsCategories
	if condition != "" {
		conditions = `(`+condition+`) AND p.`+p.IsDeleted+` = FALSE`
	} else {
		conditions = `p.`+p.IsDeleted+` = FALSE`
	}
	query = `
		SELECT p.`+p.ID+`, p.`+p.AuthorID+`, cl.`+cl.UserName+`,
				p.`+p.Title+`, p.`+p.Content+`, p.`+p.CreationDate+`,
				p.`+p.UpdateDate+`, p.`+p.DeletionDate+`,
				p.`+p.IsDeleted+`, p.`+p.Likes+`, p.`+p.Dislikes+`
		FROM `+p.Posts+` p
		JOIN `+cl.Clients+` cl ON p.`+p.AuthorID+` = cl.`+cl.ID+`
		LEFT JOIN `+pc.PostsCategories+` pc ON pc.`+pc.PostID+` = p.`+p.ID+`
		WHERE `+conditions+`;
	`
	rows, err = DB.Query(query, args...)
	return rows, err
}

func fillPostExternalData(post *Post) {
	var	err	error

	post.Comments, err = GetCommentsByPostID(post.ID)
	if err != nil {
		log.Println(err.Error())
	}
	post.Categories, err = GetPostCategories(post.ID)
	if err != nil {
		log.Println(err.Error())
	}
}

func getPostsWithCondition(condition string, args ...any) ([]*Post, error) {
	var	posts	[]*Post
	var	post	*Post
	var	rows	*sql.Rows
	var	err		error

	rows, err = getPostsWithConditionQueryResult(condition, args...)
	if err != nil {
		log.Printf("Error fetching posts: %v", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		post = &Post{}
		err = rows.Scan(
			&post.ID, &post.AuthorID, &post.UserName, &post.Title,
			&post.Content, &post.CreationDate, &post.UpdateDate,
			&post.DeletionDate, &post.IsDeleted, &post.Likes, &post.Dislikes,
		)
		if err != nil {
			log.Printf("Error scanning post: %v", err)
			continue
		}
		fillPostExternalData(post)
		posts = append(posts, post)
	}
	if err = rows.Err(); err != nil {
		log.Printf("Error during row iteration: %v", err)
		return nil, err
	}
	return posts, nil
}

func GetPostByID(id int) (*Post, error) {
	var	condition	string
	var	p			config.PostsTableKeys
	var	posts		[]*Post
	var	err			error

	p = config.TableKeys.Posts
	condition = `p.`+p.ID+` = ?`
	posts, err = getPostsWithCondition(condition, id)
	if len(posts) == 0 {
		return nil, err
	}
	return posts[0], err
}

func GetPostsByCategoryID(categoryID int) ([]*Post, error) {
	var	condition	string
	var	pc			config.PostsCategoriesTableKeys

	pc = config.TableKeys.PostsCategories
	condition = `pc.`+pc.CategoryID+` = ?`
	return getPostsWithCondition(condition, categoryID)
}

// Retrieve all the posts from database
func GetAllPosts() ([]*Post, error) {
	return getPostsWithCondition("")
}

func DeletePost(postID int) error {
	var	p	config.PostsTableKeys

	p = config.TableKeys.Posts
	query := `
		UPDATE `+p.Posts+`
		SET `+p.IsDeleted+` = 1, `+p.DeletionDate+` = CURRENT_TIMESTAMP
		WHERE `+p.ID+` = ?
	`
	result, err := DB.Exec(query, postID)
	if err != nil {
		log.Printf("Error deleting post %d: %v", postID, err)
		return fmt.Errorf("failed to delete post: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retreieve affected row: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no post found with ID %d", postID)
	}

	return nil
}

func UpdatePostLikesDislikesCount(postID int) error {
	var	query				string
	var	p					config.PostsTableKeys
	var	result				sql.Result
	var	newLikesCount		int
	var	newDislikesCount	int
	var	err					error

	p = config.TableKeys.Posts
	newLikesCount, newDislikesCount, err = GetLikeDislikeCounts(postID)
	if err != nil {
		return fmt.Errorf("failed to fetch likes and dislikes counts: %v", err)
	}
	query = `
		UPDATE `+p.Posts+`
		SET `+p.Likes+` = ?, `+p.Dislikes+` = ?
		WHERE `+p.ID+` = ?;
	`
	result, err = DB.Exec(query, newLikesCount, newDislikesCount, postID)
	if err != nil {
		return fmt.Errorf("failed to update like-dislike on post: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("no row edited: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("Post %d not found", postID)
	}
	return nil
}
