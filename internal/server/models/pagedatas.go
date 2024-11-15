package models

import "forum/internal/database"

type ViewPostPageData struct {
	Post	*database.Post
}

type HomePageData struct {
	Title		string
	Posts		[]database.Post
	IsLoggedIn	bool
	UserID		int
	UserName	string
	UserRole	string
}
