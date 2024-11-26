package models

import db "forum/internal/database"

type PostWithUserConfig struct {
	Post				*db.Post
	IsLikedByUser		bool
	IsDislikedByUser	bool
}
