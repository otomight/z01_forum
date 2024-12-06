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

func (ret ReactionElemType) String() string {
	switch (ret) {
	case ReactElemType.Post:
		return "post"
	case ReactElemType.Comment:
		return "comment"
	default:
		return ""
	}
}
