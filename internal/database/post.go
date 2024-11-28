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
	var	k		config.StructTablesKeys = config.TableKeys

	query = `
		INSERT INTO `+k.Posts.Table+` (author_id, title, category, content,
						tags, creation_date, update_date, is_deleted)
		VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 0)
	`
	result, err := DB.Exec(query, post.AuthorID, post.Title,
							post.Category, post.Content, post.Tags)
	if err != nil {
		return 0, err
	}
	postId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return postId, err
}

func GetPostByID(postID int) (*Post, error) {
	var	k	config.StructTablesKeys = config.TableKeys

	// Updated query to join Posts and Clients to get the user_name
	query := `
	SELECT p.post_id, p.author_id, c.user_name, p.title, p.category, p.tags, p.content,
			p.creation_date, p.update_date, p.deletion_date, p.is_deleted,
			p.likes, p.dislikes
	FROM `+k.Posts.Table+` p
	JOIN clients c ON p.author_id = c.user_id
	WHERE p.post_id = ? AND p.is_deleted = FALSE
	`
	post := &Post{}
	err := DB.QueryRow(query, postID).Scan(
		&post.PostID, &post.AuthorID, &post.UserName, &post.Title, &post.Category,
		&post.Tags, &post.Content, &post.CreationDate, &post.UpdateDate,
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
	var	k	config.StructTablesKeys = config.TableKeys

	query := `
		SELECT p.post_id, p.author_id, c.user_name, p.title, p.category,
			p.Tags, p.content, p.creation_date, p.update_date, p.deletion_date,
			p.is_deleted, p.likes, p.dislikes
		FROM `+k.Posts.Table+` p
		JOIN clients c ON p.author_id = c.user_id
		WHERE p.is_deleted = 0 -- Only select non deleted posts
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
		err = rows.Scan(&post.PostID, &post.AuthorID, &post.UserName,
			&post.Title, &post.Category, &post.Tags, &post.Content,
			&post.CreationDate, &post.UpdateDate, &post.DeletionDate,
			&post.IsDeleted, &post.Likes, &post.Dislikes)
		if err != nil {
			log.Printf("Error scanning post: %v", err)
			continue
		}
		post.Comments, err = GetCommentsByPostID(post.PostID)
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

// Delete post
func DeletePost(postID int) error {
	var	k	config.StructTablesKeys = config.TableKeys

	query := `
		UPDATE `+k.Posts.Table+`
		SET is_deleted = 1, deletion_date = CURRENT_TIMESTAMP
		WHERE post_id = ?
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

func UpdatePostLikesDislikesCount(postId int) error {
	var	query				string
	var	k					config.StructTablesKeys = config.TableKeys
	var	result				sql.Result
	var	newLikesCount		int
	var	newDislikesCount	int
	var	err					error

	newLikesCount, newDislikesCount, err = GetLikeDislikeCounts(postId)
	if err != nil {
		return fmt.Errorf("failed to fetch likes and dislikes counts: %v", err)
	}
	query = `
		UPDATE `+k.Posts.Table+`
		SET likes = ?, dislikes = ?
		WHERE post_id = ?;
	`
	result, err = DB.Exec(query, newLikesCount, newDislikesCount, postId)
	if err != nil {
		return fmt.Errorf("failed to update like-dislike on post: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("no row edited: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("Post %d not found", postId)
	}
	return nil
}
