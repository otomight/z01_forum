package database

import (
	"fmt"
	"forum/internal/config"
	"log"
)

func AddComment(postID, userID int, content string) error {
	var	c	config.CommentsTableKeys

	c = config.TableKeys.Comments
	query := `
		INSERT INTO `+c.Comments+` (`+c.PostID+`, `+c.UserID+`, `+c.Content+`)
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
	var	c	config.CommentsTableKeys
	var	cl	config.ClientsTableKeys

	c = config.TableKeys.Comments
	cl = config.TableKeys.Clients
	query := `
		SELECT c.`+c.ID+`, c.`+c.PostID+`, c.`+c.UserID+`,
				cl.`+cl.UserName+`, c.`+c.Content+`, c.`+c.CreationDate+`
		FROM `+c.Comments+` c
		INNER JOIN `+cl.Clients+` cl ON c.`+c.UserID+` = cl.`+cl.ID+`
		WHERE c.`+c.PostID+` = ?
		ORDER BY c.`+c.CreationDate+` ASC;
	`

	rows, err := DB.Query(query, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch comment: %w", err)
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID,
					&comment.UserName, &comment.Content, &comment.CreationDate)
		if err != nil {
			return nil, fmt.Errorf("error scanning comment: %w", err)
		}
		comments = append(comments, comment)
	}
	return comments, nil
}
