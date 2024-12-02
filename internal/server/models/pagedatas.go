package models

import db "forum/internal/database"

type ViewPostPageData struct {
	Session		*db.UserSession
	Post		*db.Post
}

type CreatePostPageData struct {
	Session		*db.UserSession
	Categories	[]*db.Category
}

type ListPostsPageData struct {
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
	Posts		[]*db.Post
}
