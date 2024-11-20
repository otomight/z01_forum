package models

type LikeDislikePostResponseAjax struct {
	LikeCount		int	`json:"like_count"`
	DislikeCount	int	`json:"dislike_count"`
}
