package database

import (
	"fmt"
	"forum/internal/config"
	"log"
)

func AddComment(postID, userID int, content string) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (post_id, user_id, content)
		VALUES(?, ?, ?)
	`, config.Table.Comments.Name)
	_, err := DB.Exec(query, postID, userID, content)
	if err != nil {
		log.Printf("Error adding comment tp post %d: %v", postID, err)
		return fmt.Errorf("failed to add comment: %w", err)
	}
	return nil
}

func GetCommentsByPostID(postID int) ([]Comment, error) {
	query := fmt.Sprintf(`
	SELECT c.comment_id, c.post_id, c.user_id, u.user_name, c.content, c.creation_date
	FROM %s c
	INNER JOIN %s u ON c.user_id = u.user_id
	WHERE c.post_id = ?
	ORDER BY c.creation_date ASC;
	`, config.Table.Comments.Name, config.Table.Clients.Name)

	rows, err := DB.Query(query, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch comment: %w", err)
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.CommentID, &comment.PostID, &comment.UserID,
					&comment.UserName, &comment.Content, &comment.CreationDate)
		if err != nil {
			return nil, fmt.Errorf("error scanning comment: %w", err)
		}
		comments = append(comments, comment)
	}
	return comments, nil
}
