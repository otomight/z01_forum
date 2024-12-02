package models

type ReactionPostResponseAjax struct {
	Added		bool	`json:"added"`
	Deleted		bool	`json:"deleted"`
	Replaced	bool	`json:"replaced"`
}
