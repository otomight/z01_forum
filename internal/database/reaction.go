package database

import (
	"database/sql"
	"fmt"
	"forum/internal/config"
	"log"
)

func AddReaction(postId int, userId int, liked bool) error {
	var	l	config.PostsReactionsTableKeys
	var	err	error

	l = config.TableKeys.PostsReactions
	_, err = insertInto(InsertIntoQuery{
		Table:		l.PostsReactions,
		Keys: []string{l.PostID, l.UserID, l.Liked},
		Values: [][]any{{
			postId, userId, liked,
		}},
		Ending: `
			ON CONFLICT(`+l.PostID+`, `+l.UserID+`) DO UPDATE
			SET `+l.Liked+` = excluded.`+l.Liked+`,
				`+l.UpdateDate+` = CURRENT_TIMESTAMP
		`,
	})
	if err != nil {
		log.Printf("Error adding like to post %d: %v", postId, err)
		return fmt.Errorf("failed to add like: %w", err)
	}
	return nil
}

func GetReactionByUser(postId int, userId int) (*Reaction, error) {
	var	query	string
	var	l		config.PostsReactionsTableKeys
	var	rows	*sql.Rows
	var	err		error
	var	ldl		Reaction

	l = config.TableKeys.PostsReactions
	query = `
		SELECT `+l.ID+`, `+l.PostID+`,
				`+l.UserID+`, `+l.Liked+`, `+l.UpdateDate+`
		FROM `+l.PostsReactions+`
		WHERE `+l.PostID+` = ? AND `+l.UserID+` = ?;
	`
	rows, err = DB.Query(query, postId, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Reaction: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&ldl.ID, &ldl.PostID,
			&ldl.UserID, &ldl.Liked, &ldl.UpdateDate)
		if err != nil {
			return nil, fmt.Errorf("error scanning Reaction: %w", err)
		}
		return &ldl, nil
	} else {
		return nil, nil
	}
}

// return like and dislike counts
func GetReactionsCounts(postId int) (int, int, error) {
	var	query			string
	var	l				config.PostsReactionsTableKeys
	var	likesCount		int
	var	dislikesCount	int
	var	err				error

	l = config.TableKeys.PostsReactions
	query = `
		SELECT
			COUNT(CASE WHEN `+l.Liked+` = 1 THEN 1 END) AS likes_count,
			COUNT(CASE WHEN `+l.Liked+` = 0 THEN 1 END) AS dislikes_count
		FROM `+l.PostsReactions+`
		WHERE `+l.PostID+` = ?;
	`
	err = DB.QueryRow(query, postId).Scan(&likesCount, &dislikesCount)
	if err != nil {
		return 0, 0, err
	}
	return likesCount, dislikesCount, nil
}

func DeleteReaction(postId int, userId int) error {
	var	query	string
	var	l		config.PostsReactionsTableKeys

	l = config.TableKeys.PostsReactions
	query = `
		DELETE FROM `+l.PostsReactions+`
		WHERE `+l.PostID+` = ? AND `+l.UserID+` = ?;
	`
	_, err := DB.Exec(query, postId, userId)
	if err != nil {
		log.Printf("Error deleting reaction of post %d: %v", postId, err)
		return fmt.Errorf("failed to delete reaction: %w", err)
	}
	return nil
}
