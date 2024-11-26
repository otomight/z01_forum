package services

import (
	db "forum/internal/database"
	"log"
)

// returns IsLikedByUser and IsDislikedByUser
func GetUserLikesConfigsOfPost(
	session	*db.UserSession,
	post	*db.Post,
) (bool, bool) {
	var	ldl		*db.LikeDislike
	var	err		error

	if session == nil {
		return false, false
	}
	ldl, err = db.GetLikeDislikeByUser(post.PostID, session.UserID)
	if err != nil {
		log.Println(err.Error())
		return false, false
	}
	if ldl == nil {
		return false, false
	}
	if ldl.Liked {
		return true, false
	}
	return false, true
}
