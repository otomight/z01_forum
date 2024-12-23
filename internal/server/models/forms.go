package models

import "forum/internal/utils"

type CreatePostForm struct {
	Title		string			`form:"title"`
	Content		string			`form:"content"`
	Categories	[]string		`form:"categories"`
	Image		*utils.FormFile	`form:"image"`
}

type RegisterForm struct {
	UserName	string	`form:"user_name"`
	Email		string	`form:"email"`
	Password	string	`form:"password"`
}

type DeletePostForm struct {
	PostID	string	`form:"post_id"`
}

type DeleteCommentForm struct {
	CommentID	string	`form:"comment_id"`
}

type LoginForm struct {
	Username	string	`form:"username"`
	Password	string	`form:"password"`
}
