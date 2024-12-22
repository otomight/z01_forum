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

func getCommentsQueryResult(
	curUserID int, condition string, args ...any,
) (*sql.Rows, error) {
	var	c			config.CommentsTableKeys
	var	cl			config.ClientsTableKeys
	var	cr			config.CommentsReactionsTableKeys

	c = config.TableKeys.Comments
	cl = config.TableKeys.Clients
	cr = config.TableKeys.CommentsReactions
	if condition != "" {
		condition = ` WHERE `+condition+``
	}
	query := `
		SELECT c.`+c.ID+`, c.`+c.PostID+`, c.`+c.UserID+`,
				cl.`+cl.UserName+`, c.`+c.Content+`, c.`+c.CreationDate+`,
				c.`+c.Likes+`, c.`+c.Dislikes+`, cr.`+cr.Liked+`
		FROM `+c.Comments+` c
		INNER JOIN `+cl.Clients+` cl ON c.`+c.UserID+` = cl.`+cl.ID+`
		LEFT JOIN `+cr.CommentsReactions+` cr
			ON cr.`+cr.CommentID+` = c.`+c.ID+` AND cr.`+cr.UserID+` = ?
		`+condition+`
		ORDER BY c.`+c.CreationDate+` ASC;
	`
	return DB.Query(query, append([]any{curUserID}, args...)...)
}

func getCommentsWithCondition(
	curUserID int, condition string, args ...any,
) ([]*Comment, error) {
	var	comments	[]*Comment
	var	comment		*Comment
	var	rows		*sql.Rows
	var	userLiked	*bool
	var	err			error

	rows, err = getCommentsQueryResult(curUserID, condition, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch comment: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		comment = &Comment{}
		err = rows.Scan(
			&comment.ID, &comment.PostID, &comment.UserID,
			&comment.UserName, &comment.Content, &comment.CreationDate,
			&comment.Likes, &comment.Dislikes, &userLiked,
		)
		if err != nil {
			log.Printf("Error scanning comment: %v\n", err)
			continue
		}
		comment.UserConfig = getUserConfig(userLiked)
		comments = append(comments, comment)
	}
	if err = rows.Err(); err != nil {
		log.Println("Error during row iteration")
		return nil, err
	}
	return comments, nil
}

func GetCommentsByPostID(curUserID int, postID int) ([]*Comment, error) {
	var	condition	string
	var	c			config.CommentsTableKeys

	c = config.TableKeys.Comments
	condition = `c.`+c.PostID+` = ?`
	return getCommentsWithCondition(curUserID, condition, postID)
}

func GetCommentsOfPostFromUser(
	curUserID int, postID int, userID int,
) ([]*Comment, error) {
	var	condition	string
	var	c			config.CommentsTableKeys

	c = config.TableKeys.Comments
	condition = `c.`+c.PostID+` = ? AND c.`+c.UserID+` = ?`
	return getCommentsWithCondition(curUserID, condition, postID, userID)
}

func deleteCommentWithCondition(condition string, args ...any) error {
	var	query	string
	var	c		config.CommentsTableKeys
	var	err		error

	c = config.TableKeys.Comments
	query = `
		DELETE FROM `+c.Comments+`
	`
	if condition != "" {
		query += ` WHERE `+condition+``
	}
	query += ";"
	_, err = DB.Exec(query, args...)
	if err != nil {
		log.Println("Error deleting comment")
		return err
	}
	return nil
}

// func DeleteComment(commentID int) {
// 	var	c			config.CommentsTableKeys
// 	var	condition	string

// 	c = config.TableKeys.Comments
// 	condition = ``+c.ID+` = ?`
// 	deleteCommentWithCondition(condition, commentID)
// }
