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
	SET liked = excluded.liked, update_date = CURRENT_TIMESTAMP
	`
	_, err := DB.Exec(query, postId, userId, liked)
	if err != nil {
		log.Printf("Error adding like to post %d: %v", postId, err)
		return fmt.Errorf("failed to add like: %w", err)
	}
	return nil
}

func GetLikeDislike(postId int, userId int) ([]LikeDislike, error) {
	var query			string
	var rows			*sql.Rows
	var err				error
	var likesDislikes	[]LikeDislike
	var ldl		LikeDislike

	query = `
	SELECT id, post_id, user_id, update_date
	FROM likes_dislikes
	WHERE post_id = ? AND user_id = ?
	`
	rows, err = DB.Query(query, postId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch LikeDislike: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		ldl = LikeDislike{}
		err = rows.Scan(&ldl.Id, &ldl.PostId, &ldl.UserId, &ldl.UpdateDate)
		if err != nil {
			return nil, fmt.Errorf("error scanning LikeDislike: %w", err)
		}
		likesDislikes = append(likesDislikes, ldl)
	}
	return likesDislikes, nil
}
