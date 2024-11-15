package models

import "forum/internal/database"

type ViewPostData struct {
	Post	*database.Post
}
