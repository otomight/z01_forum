package database

import "log"

// Post CRUD operations
func CreatePost(post *Post) error {
	_, err := DB.Exec(`
		INSERT INTO Posts (author_id, title, category, content)
		VALUES (?, ?, ?, ?)
	`, post.AuthorID, post.Title, post.Category, post.Content)
	return err
}

func GetPostByID(postID int) (*Post, error) {
	row := DB.QueryRow(`
		SELECT post_id, author_id, title, category, content, creation_date, update_date, deletion_date, is_deleted
		FROM Posts WHERE post_id = ?
	`, postID)
	var post Post
	err := row.Scan(&post.PostID, &post.AuthorID, &post.Title, &post.Category, &post.Content,
		&post.CreationDate, &post.UpdateDate, &post.DeletionDate, &post.IsDeleted)
	return &post, err
}

// Retrieve all the posts from database
func GetAllPosts() ([]Post, error) {
	query := `
		SELECT post_id, author_id, title, category, content, creation_date, update_date, deletion_date, is_deleted
		FROM Posts 
		WHERE is_deleted = 0 -- Only select non deleted posts
	`

	row, err := DB.Query(query)
	if err != nil {
		log.Printf("Error fetching posts: %v", err)
		return nil, err
	}
	defer row.Close()

	var posts []Post

	//Check iteration errors
	if err := row.Err(); err != nil {
		log.Printf("Error during row iteration: %v", err)
		return nil, err
	}

	return posts, nil
}
