package config

type ReactionElemType int

type ReactElemTypeStruct struct {
	Post	ReactionElemType
	Comment	ReactionElemType
}

const (
	post ReactionElemType = iota
	comment
)

var ReactElemType = ReactElemTypeStruct{
	Post:		post,
	Comment:	comment,
}


