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

func ViewPostHandler(w http.ResponseWriter, r *http.Request) {
	var postIdStr	string
	var postId		int
	var post		*db.Post
	var data		models.ViewPostPageData
	var session		*db.UserSession
	var err			error

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
	session, _ = r.Context().Value(config.SessionKey).(*db.UserSession)
	data = models.ViewPostPageData{
		Post:		post,
		Session:	session,
	}
	templates.RenderTemplate(w, config.ViewPostTmpl, data)
}
