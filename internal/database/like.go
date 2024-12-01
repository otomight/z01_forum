package database

import (
	"database/sql"
	"fmt"
	"forum/internal/config"
	"log"
)

func AddLikeDislike(postId int, userId int, liked bool) error {
	var	l	config.LikesDislikesTableKeys
	var	err	error

	l = config.TableKeys.LikesDislikes
	_, err = inserInto(InsertIntoQuery{
		Table:		l.LikesDislikes,
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

func GetLikeDislikeByUser(postId int, userId int) (*LikeDislike, error) {
	var	query	string
	var	l		config.LikesDislikesTableKeys
	var	rows	*sql.Rows
	var	err		error
	var	ldl		LikeDislike

	l = config.TableKeys.LikesDislikes
	query = `
		SELECT `+l.ID+`, `+l.PostID+`,
				`+l.UserID+`, `+l.Liked+`, `+l.UpdateDate+`
		FROM `+l.LikesDislikes+`
		WHERE `+l.PostID+` = ? AND `+l.UserID+` = ?;
	`
	rows, err = DB.Query(query, postId, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch LikeDislike: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&ldl.ID, &ldl.PostID,
			&ldl.UserID, &ldl.Liked, &ldl.UpdateDate)
		if err != nil {
			return nil, fmt.Errorf("error scanning LikeDislike: %w", err)
		}
		return &ldl, nil
	} else {
		return nil, nil
	}
}

func GetLikeDislikeCounts(postId int) (int, int, error) {
	var	query			string
	var	l				config.LikesDislikesTableKeys
	var	likesCount		int
	var	dislikesCount	int
	var	err				error

	l = config.TableKeys.LikesDislikes
	query = `
		SELECT
			COUNT(CASE WHEN `+l.Liked+` = 1 THEN 1 END) AS likes_count,
			COUNT(CASE WHEN `+l.Liked+` = 0 THEN 1 END) AS dislikes_count
		FROM `+l.LikesDislikes+`
		WHERE `+l.PostID+` = ?;
	`
	err = DB.QueryRow(query, postId).Scan(&likesCount, &dislikesCount)
	if err != nil {
		return 0, 0, err
	}
	return likesCount, dislikesCount, nil
}

func DeleteLikeDislike(postId int, userId int) error {
	var	query	string
	var	l		config.LikesDislikesTableKeys

	l = config.TableKeys.LikesDislikes
	query = `
		DELETE FROM `+l.LikesDislikes+`
		WHERE `+l.PostID+` = ? AND `+l.UserID+` = ?;
	`
	_, err := DB.Exec(query, postId, userId)
	if err != nil {
		log.Printf("Error deleting like-dislike of post %d: %v", postId, err)
		return fmt.Errorf("failed to delete like-dislike: %w", err)
	}
	return nil
}
