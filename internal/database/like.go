package database

import (
	"database/sql"
	"fmt"
	"log"
)

func AddLikeDislike(postId int, userId int, liked bool) error {
	query := `
	INSERT INTO likes_dislikes (post_id, user_id, liked)
	VALUES(?, ?, ?)
	ON CONFLICT(post_id, user_id) DO UPDATE
	SET liked = excluded.liked, update_date = CURRENT_TIMESTAMP;
	`
	_, err := DB.Exec(query, postId, userId, liked)
	if err != nil {
		log.Printf("Error adding like to post %d: %v", postId, err)
		return fmt.Errorf("failed to add like: %w", err)
	}
	return nil
}

func GetLikeDislikeByUser(postId int, userId int) (*LikeDislike, error) {
	var	query	string
	var	rows	*sql.Rows
	var	err		error
	var	ldl		LikeDislike

	query = `
	SELECT id, post_id, user_id, liked, update_date
	FROM likes_dislikes
	WHERE post_id = ? AND user_id = ?;
	`
	rows, err = DB.Query(query, postId, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch LikeDislike: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&ldl.Id, &ldl.PostId,
						&ldl.UserId, &ldl.Liked, &ldl.UpdateDate)
		if err != nil {
			return nil, fmt.Errorf("error scanning LikeDislike: %w", err)
		}
		return &ldl, nil
	} else {
		return nil, nil
	}
}

func	GetLikeDislikeCounts(postId int) (int, int, error) {
	var	query			string
	var	likesCount		int
	var	dislikesCount	int
	var	err				error

	query = `
		SELECT
			COUNT(CASE WHEN liked = 1 THEN 1 END) AS likes_count,
			COUNT(CASE WHEN liked = 0 THEN 1 END) AS dislikes_count
		FROM likes_dislikes
		WHERE post_id = ?;
	`
	err = DB.QueryRow(query, postId).Scan(&likesCount, &dislikesCount)
	if err != nil {
		return 0, 0, err
	}
	return likesCount, dislikesCount, nil
}

func DeleteLikeDislike(postId int, userId int) error {
	var	query	string

	query = `
	DELETE FROM likes_dislikes
	WHERE post_id = ? AND user_id = ?;
	`
	_, err := DB.Exec(query, postId, userId)
	if err != nil {
		log.Printf("Error deleting like-dislike of post %d: %v", postId, err)
		return fmt.Errorf("failed to delete like-dislike: %w", err)
	}
	return nil
}
