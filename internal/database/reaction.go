package database

import (
	"database/sql"
	"fmt"
	"forum/internal/config"
	"log"
)

func insertReactionInto(
	tableKey string, elemIDKey string, userIDKey string, likedKey string,
	updateDateKey string, elemID int, userID int, liked bool,
) error {
	var	err	error

	_, err = insertInto(InsertIntoQuery{
		Table:		tableKey,
		Keys: 		[]string{elemIDKey, userIDKey, likedKey},
		Values:		[][]any{{elemID, userID, liked}},
		Ending: `
			ON CONFLICT(`+elemIDKey+`, `+userIDKey+`) DO UPDATE
			SET `+likedKey+` = excluded.`+likedKey+`,
				`+updateDateKey+` = CURRENT_TIMESTAMP
		`,
	})
	return err
}

func AddReaction(
	elemType config.ReactionElemType,
	elemID int, userID int, liked bool,
) error {
	var	err	error
	var	pr	config.PostsReactionsTableKeys
	var	cr	config.CommentsReactionsTableKeys

	if elemType == config.ReactElemType.Post {
		pr = config.TableKeys.PostsReactions
		err = insertReactionInto(
			pr.PostsReactions, pr.PostID, pr.UserID,
			pr.Liked, pr.UpdateDate, elemID, userID, liked,
		)
	} else if elemType == config.ReactElemType.Comment {
		cr = config.TableKeys.CommentsReactions
		err = insertReactionInto(
			cr.CommentsReactions, cr.CommentID, cr.UserID,
			cr.Liked, cr.UpdateDate, elemID, userID, liked,
		)
	}
	if err != nil {
		log.Printf(
			"Error adding like to %s %d: %v", elemType.String(), elemID, err,
		)
		return fmt.Errorf("failed to add like: %w", err)
	}
	return nil
}

func getReactionByUserQuery(
	idKey string, elemIDKey string, userIDKey string,
	likedKey string, updateDateKey string, tableKey string,
) string {
	return `
		SELECT `+idKey+`, `+elemIDKey+`,
				`+userIDKey+`, `+likedKey+`, `+updateDateKey+`
		FROM `+tableKey+`
		WHERE `+elemIDKey+` = ? AND `+userIDKey+` = ?;
	`
}

func GetReactionByUser(
	elemType config.ReactionElemType, elemID int, userID int,
) (*Reaction, error) {
	var	pr		config.PostsReactionsTableKeys
	var	cr		config.CommentsReactionsTableKeys
	var	query	string
	var	rows	*sql.Rows
	var	err		error
	var	ldl		Reaction

	if elemType == config.ReactElemType.Post {
		pr = config.TableKeys.PostsReactions
		query = getReactionByUserQuery(
			pr.ID, pr.PostID, pr.UserID,
			pr.Liked, pr.UpdateDate, pr.PostsReactions,
		)
	} else if elemType == config.ReactElemType.Comment {
		cr = config.TableKeys.CommentsReactions
		query = getReactionByUserQuery(
			cr.ID, cr.CommentID, cr.UserID,
			cr.Liked, cr.UpdateDate, cr.CommentsReactions,
		)
	}
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

func getReactionCountsQuery(
	likedKey string, tableKey string, elemIDKey string,
) string {
	return `
		SELECT
			COUNT(CASE WHEN `+likedKey+` = 1 THEN 1 END) AS likes_count,
			COUNT(CASE WHEN `+likedKey+` = 0 THEN 1 END) AS dislikes_count
		FROM `+tableKey+`
		WHERE `+elemIDKey+` = ?;
	`
}

// return like and dislike counts
func getReactionsCounts(
	elemType config.ReactionElemType, elemID int,
) (int, int, error) {
	var	pr				config.PostsReactionsTableKeys
	var	cr				config.CommentsReactionsTableKeys
	var	query			string
	var	likesCount		int
	var	dislikesCount	int
	var	err				error

	if elemType == config.ReactElemType.Post {
		pr = config.TableKeys.PostsReactions
		query = getReactionCountsQuery(pr.Liked, pr.PostsReactions, pr.PostID)
	} else if elemType == config.ReactElemType.Comment {
		cr = config.TableKeys.CommentsReactions
		query = getReactionCountsQuery(
			cr.Liked, cr.CommentsReactions, cr.CommentID,
		)
	}
	err = DB.QueryRow(query, elemID).Scan(&likesCount, &dislikesCount)
	if err != nil {
		return 0, 0, err
	}
	return likesCount, dislikesCount, nil
}

func updateReactionsCountQuery(
	tableKey string, likesKey string, dislikesKey string, idKey string,
) string {
	return `
		UPDATE `+tableKey+`
		SET `+likesKey+` = ?, `+dislikesKey+` = ?
		WHERE `+idKey+` = ?;
	`
}

func UpdateReactionsCount(
	elemType config.ReactionElemType, elemID int,
) error {
	var	p					config.PostsTableKeys
	var	c					config.CommentsTableKeys
	var	query				string
	var	result				sql.Result
	var	newLikesCount		int
	var	newDislikesCount	int
	var	err					error

	if elemType == config.ReactElemType.Post {
		p = config.TableKeys.Posts
		query = updateReactionsCountQuery(p.Posts, p.Likes, p.Dislikes, p.ID)
	} else if elemType == config.ReactElemType.Comment {
		c = config.TableKeys.Comments
		query = updateReactionsCountQuery(c.Comments, c.Likes, c.Dislikes, c.ID)
	}
	newLikesCount, newDislikesCount, err = getReactionsCounts(elemType, elemID)
	if err != nil {
		return fmt.Errorf("failed to fetch likes and dislikes counts: %v", err)
	}
	result, err = DB.Exec(query, newLikesCount, newDislikesCount, elemID)
	if err != nil {
		return fmt.Errorf(
			"failed to update reactions on %s: %w", elemType.String(), err,
		)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("no row edited: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s %d not found", elemType.String(), elemID)
	}
	return nil
}

func deleteReactionsWithCondition(
	tableKey string, condition string, args ...any,
) error {
	var	query	string
	var	err		error

	query = `
		DELETE FROM `+tableKey+`
	`
	if condition != "" {
		query += ` WHERE `+condition+``
	}
	query += ";"
	_, err = DB.Exec(query, args...)
	if err != nil {
		log.Printf("Error deleting reaction of table %s\n", tableKey)
		return err
	}
	return nil
}

func DeleteUserReaction(
	elemType config.ReactionElemType, elemID int, userID int,
) error {
	var	pr			config.PostsReactionsTableKeys
	var	cr			config.CommentsReactionsTableKeys
	var	tableKey	string
	var	condition	string

	if elemType == config.ReactElemType.Post {
		pr = config.TableKeys.PostsReactions
		tableKey = pr.PostsReactions
		condition = ``+pr.PostID+` = ? AND `+pr.UserID+` = ?`
	} else if elemType == config.ReactElemType.Comment {
		cr = config.TableKeys.CommentsReactions
		tableKey = cr.CommentsReactions
		condition = ``+cr.CommentID+` = ? AND `+cr.UserID+` = ?`
	}
	return deleteReactionsWithCondition(tableKey, condition, elemID, userID)
}
