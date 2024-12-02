package posthandlers

import (
	"forum/internal/config"
	db "forum/internal/database"
	"forum/internal/server/models"
	"forum/internal/server/templates"
	"net/http"
	"strconv"
	"strings"
)

// func fillViewPostPageData(
// 	post	*db.Post,
// 	ctx		context.Context,
// ) *models.ViewPostPageData {
// 	var	session				*db.UserSession
// 	var	data				models.ViewPostPageData
// 	var	postWithUserConfig	models.PostWithUserConfig
// 	var	isLikedByUser		bool
// 	var	isDislikedByUser	bool

// 	session, _ = ctx.Value(config.SessionKey).(*db.UserSession)
// 	isLikedByUser, isDislikedByUser =
// 			services.GetUserLikesConfigsOfPost(session, post)
// 	postWithUserConfig = models.PostWithUserConfig{
// 		Post:				post,
// 		IsLikedByUser:		isLikedByUser,
// 		IsDislikedByUser:	isDislikedByUser,
// 	}
// 	data = models.ViewPostPageData{
// 		PostWithUserConfig:	&postWithUserConfig,
// 		Session:			session,
// 	}
// 	return &data
// }

func ViewPostHandler(w http.ResponseWriter, r *http.Request) {
	var	postIDStr			string
	var	postID				int
	var	post				*db.Post
	var	data				*models.ViewPostPageData
	var	session				*db.UserSession
	var	userID				int
	var	err					error

	if r.Method != http.MethodGet {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}
	postIDStr = strings.TrimPrefix(r.URL.Path, "/post/view/")
	if postIDStr == "" || strings.Contains(postIDStr, "/") {
		http.NotFound(w, r)
		return
	}
	if postID, err = strconv.Atoi(postIDStr); err != nil {
		http.NotFound(w, r)
		return
	}
	session, _ = r.Context().Value(config.SessionKey).(*db.UserSession)
	if session == nil {
		userID = 0
	} else {
		userID = session.UserID
	}
	if post, err = db.GetPostByID(userID, postID); err != nil {
		http.NotFound(w, r)
		return
	}
	data = &models.ViewPostPageData{
		Session:	session,
		Post:		post,
	}
	templates.RenderTemplate(w, config.ViewPostTmpl, data)
}
