package models

type CreatePostForm struct {
	Title		string	`form:"title"`
	Content		string	`form:"content"`
	Category	string	`form:"category"`
	Tags		string	`form:"tags"`
}

type RegisterForm struct {
	UserName	string	`form:"user_name"`
	Email		string	`form:"email"`
	Password	string	`form:"password"`
	FirstName	string	`form:"first_name"`
	LastName	string	`form:"last_name"`
}

type DeletePostForm struct {
	PostId	string	`form:"postId"`
}

type LoginForm struct {
	Username	string	`form:"username"`
	Password	string	`form:"password"`
}
