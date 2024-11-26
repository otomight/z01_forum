package models

import "forum/internal/database"

type ViewPostPageData struct {
	Session				*database.UserSession
	Post				*database.Post
	IsLikedByUser		bool
	IsDislikedByUser	bool
}

type CreatePostPageData struct {
	Session	*database.UserSession
}

type HomePageData struct {
	Posts		[]*PostWithUserConfig
	Session		*database.UserSession
}
