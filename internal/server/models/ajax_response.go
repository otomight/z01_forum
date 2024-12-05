package models

type ReactionResponseAjax struct {
	Added		bool	`json:"added"`
	Deleted		bool	`json:"deleted"`
	Replaced	bool	`json:"replaced"`
}
