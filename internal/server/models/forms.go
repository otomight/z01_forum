package models

type CreationPostForm struct {
	Title		string	`form:"title"`
	Content		string	`form:"content"`
	Category	string	`form:"category"`
	Tags		string	`form:"tags"`
}
