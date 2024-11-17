package database

import (
	"fmt"
	"log"
)

func AddComment(postID, userID int, content string) error {
	query := `
	INSERT INTO Comments (post_id, user_id, content)
	VALUES(?, ?, ?)
	`
	_, err := DB.Exec(query, postID, userID, content)
	if err != nil {
		log.Printf("Error adding comment tp post %d: %v", postID, err)
		return fmt.Errorf("failed to add comment: %w", err)
	}
	return nil
}

func GetCommentsByPostID(postID int) ([]Comment, error) {
	query := `
	SELECT comment_id, post_id, user_id, user_name, content, creation_date
	FROM Comment
	WHERE post_id = ? AND is_deleted = FALSE
	ORDER BY creation_date ASC
	`

	rows, err := DB.Query(query, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch comment: %w", err)
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.CommentID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreationDate, &comment.UserName); err != nil {
			return nil, fmt.Errorf("error scanning comment: %w", err)
		}
		comments = append(comments, comment)
	}
	return comments, nil
}
