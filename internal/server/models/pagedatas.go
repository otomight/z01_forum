package models

import db "forum/internal/database"

type ViewPostPageData struct {
	Session				*db.UserSession
	PostWithUserConfig	*PostWithUserConfig
}

type CreatePostPageData struct {
	Session		*db.UserSession
	Categories	[]*db.Category
}

type PostsInCategoriesPageData struct {
	Session		*db.UserSession
	Category	*db.Category
	Posts		[]*db.Post
}

type CategoriesPageData struct {
	Session		*db.UserSession
	Categories	[]*db.Category
}

type HomePageData struct {
	Session		*db.UserSession
	Posts		[]*PostWithUserConfig
}
