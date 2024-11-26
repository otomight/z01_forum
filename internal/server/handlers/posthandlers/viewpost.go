package posthandlers

import (
	"context"
	"forum/internal/config"
	db "forum/internal/database"
	"forum/internal/server/models"
	"forum/internal/server/services"
	"forum/internal/server/templates"
	"net/http"
	"strconv"
	"strings"
)

func fillViewPostPageData(
	post	*db.Post,
	ctx		context.Context,
) *models.ViewPostPageData {
	var	session				*db.UserSession
	var	data				models.ViewPostPageData
	var	postWithUserConfig	models.PostWithUserConfig
	var	isLikedByUser		bool
	var	isDislikedByUser	bool

	session, _ = ctx.Value(config.SessionKey).(*db.UserSession)
	isLikedByUser, isDislikedByUser =
			services.GetUserLikesConfigsOfPost(session, post)
	postWithUserConfig = models.PostWithUserConfig{
		Post:				post,
		IsLikedByUser:		isLikedByUser,
		IsDislikedByUser:	isDislikedByUser,
	}
	data = models.ViewPostPageData{
		PostWithUserConfig:	&postWithUserConfig,
		Session:			session,
	}
	return &data
}

func ViewPostHandler(w http.ResponseWriter, r *http.Request) {
	var	postIdStr			string
	var	postId				int
	var	post				*db.Post
	var	data				*models.ViewPostPageData
	var	err					error

	if r.Method != http.MethodGet {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}
	postIdStr = strings.TrimPrefix(r.URL.Path, "/post/view/")
	if postIdStr == "" || strings.Contains(postIdStr, "/") {
		http.NotFound(w, r)
		return
	}
	if postId, err = strconv.Atoi(postIdStr); err != nil {
		http.NotFound(w, r)
		return
	}
	if post, err = db.GetPostByID(postId); err != nil {
		http.NotFound(w, r)
		return
	}
	data = fillViewPostPageData(post, r.Context())
	templates.RenderTemplate(w, config.ViewPostTmpl, data)
}
