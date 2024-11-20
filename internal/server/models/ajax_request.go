package models

type LikeDislikePostRequestAjax struct {
	PostId	int	`json:"post_id"`
	UserId	int	`json:"user_id"`
}
