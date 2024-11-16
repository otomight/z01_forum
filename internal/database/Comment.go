package database

import (
	"fmt"
	"log"
)

func AddComment(postID, userID int, content string) error {
	query := `
	INSERT INTO comments (post_id, user_id, content, creation_date)
	VALUES(?, ?, ?, CURRENT_TIMESTAMP)
	`
	_, err := DB.Exec(query, postID, userID, content)
	if err != nil {
		log.Printf("Error adding comment tp post %d: %v", postID, err)
		return fmt.Errorf("failed to add comment: %w", err)
	}
	return nil
}
