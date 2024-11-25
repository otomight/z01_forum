package models

type LikeDislikePostResponseAjax struct {
	Added		bool	`json:"added"`
	Deleted		bool	`json:"deleted"`
	Replaced	bool	`json:"replaced"`
}
