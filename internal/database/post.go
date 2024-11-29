package database

import (
	"database/sql"
	"fmt"
	"forum/internal/config"
	"log"
)

// Post CRUD operations
func NewPost(post *Post) (int64, error) {
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
	return postID, err
}

func GetPostByID(postID int) (*Post, error) {
	var	p	config.PostsTableKeys
	var	c	config.ClientsTableKeys

	p = config.TableKeys.Posts
	c = config.TableKeys.Clients
	// Updated query to join Posts and Clients to get the user_name
	query := `
	SELECT p.`+p.ID+`, p.`+p.AuthorID+`, c.`+c.UserName+`,
			p.`+p.Title+`, p.`+p.Content+`, p.`+p.CreationDate+`,
			p.`+p.UpdateDate+`, p.`+p.DeletionDate+`,
			p.`+p.IsDeleted+`, p.`+p.Likes+`, p.`+p.Dislikes+`
	FROM `+p.Posts+` p
	JOIN `+c.Clients+` c ON p.`+p.AuthorID+` = c.`+c.ID+`
	WHERE p.`+p.ID+` = ? AND p.`+p.IsDeleted+` = FALSE
	`
	post := &Post{}
	err := DB.QueryRow(query, postID).Scan(
		&post.ID, &post.AuthorID, &post.UserName, &post.Title,
		&post.Content, &post.CreationDate, &post.UpdateDate,
		&post.DeletionDate, &post.IsDeleted, &post.Likes, &post.Dislikes,
	)
	if err != nil {
		log.Printf("Error retrieving post by ID %d: %v", postID, err)
		return nil, fmt.Errorf("could not fetch post: %w", err)
	}

	// Fetch associated comments for the post
	post.Comments, err = GetCommentsByPostID(postID)
	if err != nil {
		log.Println(err.Error())
	}

	return post, nil
}

// Retrieve all the posts from database
func GetAllPosts() ([]Post, error) {
	var	p	config.PostsTableKeys
	var c	config.ClientsTableKeys

	p = config.TableKeys.Posts
	c = config.TableKeys.Clients
	query := `
		SELECT p.`+p.ID+`, p.`+p.AuthorID+`, c.`+c.UserName+`,
				p.`+p.Title+`, p.`+p.Content+`, p.`+p.CreationDate+`,
				p.`+p.UpdateDate+`, p.`+p.DeletionDate+`,
				p.`+p.IsDeleted+`, p.`+p.Likes+`, p.`+p.Dislikes+`
		FROM `+p.Posts+` p
		JOIN `+c.Clients+` c ON p.`+p.AuthorID+` = c.`+c.ID+`
		WHERE p.`+p.IsDeleted+` = 0 -- Only select non deleted posts
	`
	rows, err := DB.Query(query)
	if err != nil {
		log.Printf("Error fetching posts: %v", err)
		return nil, err
	}
	defer rows.Close()
	var posts []Post
	// Iterate through rows and append to posts slice
	for rows.Next() {
		var post Post
		err = rows.Scan(
			&post.ID, &post.AuthorID, &post.UserName, &post.Title,
			&post.Content, &post.CreationDate, &post.UpdateDate,
			&post.DeletionDate, &post.IsDeleted, &post.Likes, &post.Dislikes,
		)
		if err != nil {
			log.Printf("Error scanning post: %v", err)
			continue
		}
		post.Comments, err = GetCommentsByPostID(post.ID)
		if err != nil {
			log.Println(err.Error())
		}
		posts = append(posts, post)
	}
	//Check iteration errors
	if err = rows.Err(); err != nil {
		log.Printf("Error during row iteration: %v", err)
		return nil, err
	}
	return posts, nil
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
