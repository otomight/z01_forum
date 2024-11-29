package models

type LikeDislikePostRequestAjax struct {
	PostID	int	`json:"post_id"`
	UserID	int	`json:"user_id"`
}
