package database

import (
	"fmt"
	"log"
)

// Post CRUD operations
func NewPost(post *Post) (int64, error) {
	result, err := DB.Exec(`
		INSERT INTO Posts (author_id, title, category, content, tags, creation_date, update_date, is_deleted)
		VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 0)
	`, post.AuthorID, post.Title, post.Category, post.Content, post.Tags)
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
	row := DB.QueryRow(`
		SELECT post_id, author_id, title, category, tags, content, creation_date, update_date, deletion_date, is_deleted
		FROM Posts WHERE post_id = ?
	`, postID)
	var post Post
	err := row.Scan(&post.PostID, &post.AuthorID, &post.Title, &post.Category, &post.Tags, &post.Content,
		&post.CreationDate, &post.UpdateDate, &post.DeletionDate, &post.IsDeleted)
	if err != nil {
		log.Printf("Error retrieving post by ID %d: %v", postID, err)
		return nil, err
	}
	return &post, err
}

// Retrieve all the posts from database
func GetAllPosts() ([]Post, error) {
	query := `
		SELECT Posts.post_id, Posts.author_id, Clients.user_name, Posts.title, Posts.category, Posts.Tags, Posts.content, Posts.creation_date, Posts.update_date, Posts.deletion_date, Posts.is_deleted
		FROM Posts
		JOIN Clients ON Posts.author_id = Clients.user_id
		WHERE Posts.is_deleted = 0 -- Only select non deleted posts
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
		if err := rows.Scan(&post.PostID, &post.AuthorID, &post.UserName, &post.Title, &post.Category, &post.Tags, &post.Content, &post.CreationDate, &post.UpdateDate, &post.DeletionDate, &post.IsDeleted); err != nil {
			log.Printf("Error scanning post: %v", err)
			return nil, err
		}
		posts = append(posts, post)
	}

	//Check iteration errors
	if err := rows.Err(); err != nil {
		log.Printf("Error during row iteration: %v", err)
		return nil, err
	}

	return posts, nil
}

// Delete post
func DeletePost(postID int) error {
	query := `
		UPDATE posts
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

// inserting post for testing
func InsertSamplePost() {
	// Check if any posts already exist
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM Posts").Scan(&count)
	if err != nil {
		log.Printf("Failed to count posts: %v", err)
		return
	}

	// If no posts exist, insert the sample post
	if count == 0 {
		_, err := DB.Exec(`
            INSERT INTO Posts (author_id, title, category, tags, content, creation_date, update_date, is_deleted)
            VALUES (1, 'Sample Post', 'General', 'other', 'This is a sample post content.', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 0)
        `)
		if err != nil {
			log.Printf("Failed to insert sample post: %v", err)
		}
	}
}
