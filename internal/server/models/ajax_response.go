package models

type ReactionResponseAjax struct {
	LikesCount		int		`json:"likes_count"`
	DislikesCount	int		`json:"dislikes_count"`
	Added			bool	`json:"added"`
	Deleted			bool	`json:"deleted"`
	Replaced		bool	`json:"replaced"`

}
