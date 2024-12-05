package database

import (
	"database/sql"
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

func GetCommentsByPostID(curUserID int, postID int) ([]Comment, error) {
	var	c			config.CommentsTableKeys
	var	cl			config.ClientsTableKeys
	var	cr			config.CommentsReactionsTableKeys
	var	userLiked	*bool

	c = config.TableKeys.Comments
	cl = config.TableKeys.Clients
	cr = config.TableKeys.CommentsReactions
	query := `
		SELECT c.`+c.ID+`, c.`+c.PostID+`, c.`+c.UserID+`,
				cl.`+cl.UserName+`, c.`+c.Content+`, c.`+c.CreationDate+`,
				c.`+c.Likes+`, c.`+c.Dislikes+`, cr.`+cr.Liked+`
		FROM `+c.Comments+` c
		INNER JOIN `+cl.Clients+` cl ON c.`+c.UserID+` = cl.`+cl.ID+`
		LEFT JOIN `+cr.CommentsReactions+` cr
			ON cr.`+cr.CommentID+` = c.`+c.ID+` AND cr.`+cr.UserID+` = ?
		WHERE c.`+c.PostID+` = ?
		ORDER BY c.`+c.CreationDate+` ASC;
	`

	rows, err := DB.Query(query, curUserID, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch comment: %w", err)
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(
			&comment.ID, &comment.PostID, &comment.UserID,
			&comment.UserName, &comment.Content, &comment.CreationDate,
			&comment.Likes, &comment.Dislikes, &userLiked,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning comment: %w", err)
		}
		comment.UserConfig = getUserConfig(userLiked)
		comments = append(comments, comment)
	}
	return comments, nil
}

func UpdateCommentReactionsCount(commentID int) error {
	var	query				string
	var	c					config.CommentsTableKeys
	var	result				sql.Result
	var	newLikesCount		int
	var	newDislikesCount	int
	var	err					error

	c = config.TableKeys.Comments
	newLikesCount, newDislikesCount, err =
			GetReactionsCounts(config.ReactElemType.Comment, commentID)
	if err != nil {
		return fmt.Errorf("failed to fetch likes and dislikes counts: %v", err)
	}
	query = `
		UPDATE `+c.Comments+`
		SET `+c.Likes+` = ?, `+c.Dislikes+` = ?
		WHERE `+c.ID+` = ?;
	`
	result, err = DB.Exec(query, newLikesCount, newDislikesCount, commentID)
	if err != nil {
		return fmt.Errorf("failed to update like-dislike on comment: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("no row edited: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("Post %d not found", commentID)
	}
	return nil
}
