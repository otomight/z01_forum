package models

import "forum/internal/database"

type ViewPostPageData struct {
	Session				*database.UserSession
	PostWithUserConfig	*PostWithUserConfig
}

type CreatePostPageData struct {
	Session	*database.UserSession
}

type HomePageData struct {
	Session		*database.UserSession
	Posts		[]*PostWithUserConfig
}
