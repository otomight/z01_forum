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

func fillViewPostPageData(
	w http.ResponseWriter, r *http.Request, postID int,
) (*models.ViewPostPageData, error) {
	var	session		*db.UserSession
	var	categories	[]*db.Category
	var	data		*models.ViewPostPageData
	var	userID		int
	var	post		*db.Post
	var	err			error

	session, _ = r.Context().Value(config.SessionKey).(*db.UserSession)
	if session == nil {
		userID = 0
	} else {
		userID = session.UserID
	}
	if categories, err = db.GetGlobalCategories(); err != nil {
		http.Error(
			w, "Error at fetching categories", http.StatusInternalServerError,
		)
		return nil, err
	}
	if post, err = db.GetPostByID(userID, postID); err != nil {
		http.NotFound(w, r)
		return nil, err
	}
	data = &models.ViewPostPageData{
		Session:	session,
		Categories:	categories,
		Post:		post,
	}
	return data, nil
}

func ViewPostHandler(w http.ResponseWriter, r *http.Request) {
	var	postIDStr			string
	var	postID				int
	var	data				*models.ViewPostPageData
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
	if data, err = fillViewPostPageData(w, r, postID); err != nil {
		return
	}
	templates.RenderTemplate(w, config.ViewPostTmpl, data)
}
