package database

import (
	"database/sql"
	"fmt"
	"forum/internal/config"
	"log"
)

func addReactionToPost(elemID int, userID int, liked bool) error {
	var	pr			config.PostsReactionsTableKeys
	var	err			error

	pr = config.TableKeys.PostsReactions
	_, err = insertInto(InsertIntoQuery{
		Table:		pr.PostsReactions,
		Keys: 		[]string{pr.PostID, pr.UserID, pr.Liked},
		Values:		[][]any{{elemID, userID, liked}},
		Ending: `
			ON CONFLICT(`+pr.PostID+`, `+pr.UserID+`) DO UPDATE
			SET `+pr.Liked+` = excluded.`+pr.Liked+`,
				`+pr.UpdateDate+` = CURRENT_TIMESTAMP
		`,
	})
	if err != nil {
		log.Printf("Error adding like to post %d: %v", elemID, err)
		return fmt.Errorf("failed to add like: %w", err)
	}
	return nil
}

func addReactionToComment(elemID int, userID int, liked bool) error {
	var	cr			config.CommentsReactionsTableKeys
	var	err			error

	cr = config.TableKeys.CommentsReactions
	_, err = insertInto(InsertIntoQuery{
		Table:		cr.CommentsReactions,
		Keys: 		[]string{cr.CommentID, cr.UserID, cr.Liked},
		Values:		[][]any{{elemID, userID, liked}},
		Ending: `
			ON CONFLICT(`+cr.CommentID+`, `+cr.UserID+`) DO UPDATE
			SET `+cr.Liked+` = excluded.`+cr.Liked+`,
				`+cr.UpdateDate+` = CURRENT_TIMESTAMP
		`,
	})
	if err != nil {
		log.Printf("Error adding like to comment %d: %v", elemID, err)
		return fmt.Errorf("failed to add like: %w", err)
	}
	return nil
}

func AddReaction(
	elemType config.ReactionElemType,
	elemID int, userID int, liked bool,
) error {
	if elemType == config.ReactElemType.Post {
		return addReactionToPost(elemID, userID, liked)
	} else if elemType == config.ReactElemType.Comment {
		return addReactionToComment(elemID, userID, liked)
	}
	return fmt.Errorf("Unexpected error")
}

func getReactionByUserQuery(elemType config.ReactionElemType) string {
	var	query	string
	var	pr		config.PostsReactionsTableKeys
	var	cr		config.CommentsReactionsTableKeys

	if elemType == config.ReactElemType.Post {
		pr = config.TableKeys.PostsReactions
		query = `
			SELECT `+pr.ID+`, `+pr.PostID+`,
					`+pr.UserID+`, `+pr.Liked+`, `+pr.UpdateDate+`
			FROM `+pr.PostsReactions+`
			WHERE `+pr.PostID+` = ? AND `+pr.UserID+` = ?;
		`
	} else if elemType == config.ReactElemType.Comment {
		cr = config.TableKeys.CommentsReactions
		query = `
			SELECT `+cr.ID+`, `+cr.CommentID+`,
					`+cr.UserID+`, `+cr.Liked+`, `+cr.UpdateDate+`
			FROM `+cr.CommentsReactions+`
			WHERE `+cr.CommentID+` = ? AND `+cr.UserID+` = ?;
		`
	}
	return query
}

func GetReactionByUser(
	elemType config.ReactionElemType,
	elemID int, userID int,
) (*Reaction, error) {
	var	query	string
	var	rows	*sql.Rows
	var	err		error
	var	ldl		Reaction

	query = getReactionByUserQuery(elemType)
	rows, err = DB.Query(query, elemID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Reaction: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&ldl.ID, &ldl.ElemID,
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
func GetReactionsCounts(
	elemType config.ReactionElemType, elemID int,
) (int, int, error) {
	var	query			string
	var	pr				config.PostsReactionsTableKeys
	var	cr				config.CommentsReactionsTableKeys
	var	likesCount		int
	var	dislikesCount	int
	var	err				error

	if elemType == config.ReactElemType.Post {
		pr = config.TableKeys.PostsReactions
		query = `
			SELECT
				COUNT(CASE WHEN `+pr.Liked+` = 1 THEN 1 END) AS likes_count,
				COUNT(CASE WHEN `+pr.Liked+` = 0 THEN 1 END) AS dislikes_count
			FROM `+pr.PostsReactions+`
			WHERE `+pr.PostID+` = ?;
		`
	} else if elemType == config.ReactElemType.Comment {
		cr = config.TableKeys.CommentsReactions
		query = `
			SELECT
				COUNT(CASE WHEN `+cr.Liked+` = 1 THEN 1 END) AS likes_count,
				COUNT(CASE WHEN `+cr.Liked+` = 0 THEN 1 END) AS dislikes_count
			FROM `+cr.CommentsReactions+`
			WHERE `+cr.CommentID+` = ?;
		`
	}
	err = DB.QueryRow(query, elemID).Scan(&likesCount, &dislikesCount)
	if err != nil {
		return 0, 0, err
	}
	return likesCount, dislikesCount, nil
}

func DeleteReaction(
	elemType config.ReactionElemType,
	postId int, userId int,
) error {
	var	query	string
	var	pr		config.PostsReactionsTableKeys
	var	cr		config.CommentsReactionsTableKeys

	if elemType == config.ReactElemType.Post {
		pr = config.TableKeys.PostsReactions
		query = `
			DELETE FROM `+pr.PostsReactions+`
			WHERE `+pr.PostID+` = ? AND `+pr.UserID+` = ?;
		`
	} else if elemType == config.ReactElemType.Comment {
		cr = config.TableKeys.CommentsReactions
		query = `
			DELETE FROM `+cr.CommentsReactions+`
			WHERE `+cr.CommentID+` = ? AND `+cr.UserID+` = ?;
		`
	}
	_, err := DB.Exec(query, postId, userId)
	if err != nil {
		log.Printf("Error deleting reaction of post %d: %v", postId, err)
		return fmt.Errorf("failed to delete reaction: %w", err)
	}
	return nil
}
