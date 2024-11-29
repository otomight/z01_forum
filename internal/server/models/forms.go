package models

type CreatePostForm struct {
	Title		string		`form:"title"`
	Content		string		`form:"content"`
	Categories	[]string	`form:"categories"`
}

type RegisterForm struct {
	UserName	string	`form:"user_name"`
	Email		string	`form:"email"`
	Password	string	`form:"password"`
	FirstName	string	`form:"first_name"`
	LastName	string	`form:"last_name"`
}

type DeletePostForm struct {
	PostID	string	`form:"post_id"`
}

type LoginForm struct {
	Username	string	`form:"username"`
	Password	string	`form:"password"`
}
