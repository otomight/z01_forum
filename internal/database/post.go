package database

import (
	"database/sql"
	"fmt"
	"forum/internal/config"
	"log"
)

// Post CRUD operations
func NewPost(post *Post, categoriesIDs []int) (int, error) {
	var	p		config.PostsTableKeys

	p = config.TableKeys.Posts
	result, err := insertInto(InsertIntoQuery{
		Table:		p.Posts,
		Keys: []string{
			p.AuthorID, p.Title, p.Content, p.IsDeleted,
		},
		Values: [][]any{{
			post.AuthorID, post.Title, post.Content, 0,
		}},
	})
	if err != nil {
		return 0, err
	}
	postID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	if err = AddPostCategories(int(postID), categoriesIDs...); err != nil {
		log.Printf("Error adding categories to the database: %v", err)
	}
	return int(postID), err
}

func getPostsWithConditionQueryResult(
	userID int, condition string, args ...any,
) (*sql.Rows, error) {
	var	p			config.PostsTableKeys
	var	cl			config.ClientsTableKeys
	var	pc			config.PostsCategoriesTableKeys
	var	r			config.ReactionsTableKeys
	var	query		string
	var	rows		*sql.Rows
	var	conditions	string
	var err			error

	p = config.TableKeys.Posts
	cl = config.TableKeys.Clients
	pc = config.TableKeys.PostsCategories
	r = config.TableKeys.Reactions
	if condition != "" {
		conditions = `(`+condition+`) AND p.`+p.IsDeleted+` = FALSE`
	} else {
		conditions = `p.`+p.IsDeleted+` = FALSE`
	}
	query = `
		SELECT DISTINCT p.`+p.ID+`, p.`+p.AuthorID+`, cl.`+cl.UserName+`,
			p.`+p.Title+`, p.`+p.Content+`, p.`+p.CreationDate+`,
			p.`+p.UpdateDate+`, p.`+p.DeletionDate+`, p.`+p.IsDeleted+`,
			p.`+p.Likes+`, p.`+p.Dislikes+`, r.`+r.Liked+`
		FROM `+p.Posts+` p
		JOIN `+cl.Clients+` cl ON p.`+p.AuthorID+` = cl.`+cl.ID+`
		LEFT JOIN `+pc.PostsCategories+` pc ON pc.`+pc.PostID+` = p.`+p.ID+`
		LEFT JOIN `+r.Reactions+` r
		ON r.`+r.PostID+` = p.`+p.ID+` AND r.`+r.UserID+` = ?
		WHERE `+conditions+`;
	`
	rows, err = DB.Query(query, append([]any{userID}, args...)...)
	return rows, err
}

func getUserConfig(userLiked *bool) *UserConfig {
	var	userConfig	UserConfig

	if userLiked == nil {
		return &userConfig
	}
	if *userLiked == true {
		userConfig.IsLiked = true
		return &userConfig
	} else {
		userConfig.IsDisliked = true
		return &userConfig
	}
}

func fillPostExternalData(post *Post, userLiked *bool) {
	var	err	error

	post.UserConfig = getUserConfig(userLiked)
	post.Comments, err = GetCommentsByPostID(post.ID)
	if err != nil {
		log.Println(err.Error())
	}
	post.Categories, err = GetPostCategories(post.ID)
	if err != nil {
		log.Println(err.Error())
	}
}

func getPostsWithCondition(userID int, condition string, args ...any) ([]*Post, error) {
	var	posts		[]*Post
	var	post		*Post
	var	rows		*sql.Rows
	var	userLiked	*bool
	var	err			error

	rows, err = getPostsWithConditionQueryResult(userID, condition, args...)
	if err != nil {
		log.Println("Error on post query")
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		post = &Post{}
		err = rows.Scan(
			&post.ID, &post.AuthorID, &post.UserName, &post.Title,
			&post.Content, &post.CreationDate, &post.UpdateDate,
			&post.DeletionDate, &post.IsDeleted,
			&post.Likes, &post.Dislikes, &userLiked,
		)
		if err != nil {
			log.Println("Error scanning post")
			continue
		}
		fillPostExternalData(post, userLiked)
		posts = append(posts, post)
	}
	if err = rows.Err(); err != nil {
		log.Println("Error during row iteration")
		return nil, err
	}
	return posts, nil
}

func GetPostByID(userID int, id int) (*Post, error) {
	var	condition	string
	var	p			config.PostsTableKeys
	var	posts		[]*Post
	var	err			error

	p = config.TableKeys.Posts
	condition = `p.`+p.ID+` = ?`
	posts, err = getPostsWithCondition(userID, condition, id)
	if len(posts) == 0 {
		return nil, err
	}
	return posts[0], err
}

func GetPostsRelatedToCurUser(userID int) ([]*Post, error) {
	var	condition	string
	var	p			config.PostsTableKeys
	var	r			config.ReactionsTableKeys

	p = config.TableKeys.Posts
	r = config.TableKeys.Reactions
	condition = `p.`+p.AuthorID+` = ? OR `+r.UserID+` = ?`
	return getPostsWithCondition(userID, condition, userID, userID)
}

func GetPostsByCategoryID(userID int, categoryID int) ([]*Post, error) {
	var	condition	string
	var	pc			config.PostsCategoriesTableKeys

	pc = config.TableKeys.PostsCategories
	condition = `pc.`+pc.CategoryID+` = ?`
	return getPostsWithCondition(userID, condition, categoryID)
}

// Retrieve all the posts from database
func GetAllPosts(userID int) ([]*Post, error) {
	return getPostsWithCondition(userID, "")
}

func DeletePost(postID int) error {
	var	p	config.PostsTableKeys

	p = config.TableKeys.Posts
	query := `
		UPDATE `+p.Posts+`
		SET `+p.IsDeleted+` = 1, `+p.DeletionDate+` = CURRENT_TIMESTAMP
		WHERE `+p.ID+` = ?
	`
	result, err := DB.Exec(query, postID)
	if err != nil {
		log.Printf("Error deleting post %d: %v", postID, err)
		return fmt.Errorf("failed to delete post: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retreieve affected row: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no post found with ID %d", postID)
	}

	return nil
}

func UpdatePostReactionsCount(postID int) error {
	var	query				string
	var	p					config.PostsTableKeys
	var	result				sql.Result
	var	newLikesCount		int
	var	newDislikesCount	int
	var	err					error

	p = config.TableKeys.Posts
	newLikesCount, newDislikesCount, err = GetReactionsCounts(postID)
	if err != nil {
		return fmt.Errorf("failed to fetch likes and dislikes counts: %v", err)
	}
	query = `
		UPDATE `+p.Posts+`
		SET `+p.Likes+` = ?, `+p.Dislikes+` = ?
		WHERE `+p.ID+` = ?;
	`
	result, err = DB.Exec(query, newLikesCount, newDislikesCount, postID)
	if err != nil {
		return fmt.Errorf("failed to update like-dislike on post: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("no row edited: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("Post %d not found", postID)
	}
	return nil
}
