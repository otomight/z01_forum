package handlers

import (
	"forum/internal/config"
	db "forum/internal/database"
	"forum/internal/server/models"
	"forum/internal/server/templates"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func fillListPostsPageData(
	w http.ResponseWriter, r *http.Request, categoryID int,
) (*models.CategoryPostsPageData, error) {
	var	session		*db.UserSession
	var	categories	[]*db.Category
	var	data		*models.CategoryPostsPageData
	var	userID		int
	var	category	*db.Category
	var	posts		[]*db.Post
	var	err			error

	category, err = db.GetGlobalCategoryByID(categoryID)
	if err != nil {
		return nil, err
	}
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
	if posts, err = db.GetPostsByCategoryID(userID, categoryID); err != nil {
		http.Error(
			w, "Failed to fetch posts from category",
			http.StatusInternalServerError,
		)
		return nil, err
	}
	data = &models.CategoryPostsPageData{
		Session:	session,
		Categories:	categories,
		Category:	category,
		Posts:		posts,
	}
	return data, nil
}

func CategoryPostsPageHandler(w http.ResponseWriter, r *http.Request) {
	var	idStr	string
	var	id		int
	var	err		error
	var	data	*models.CategoryPostsPageData

	if r.Method != http.MethodGet {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}
	idStr = strings.TrimPrefix(r.URL.Path, "/categories/")
	if idStr == "" || strings.Contains(idStr, "/") {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}
	id, err = strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}
	if data, err = fillListPostsPageData(w, r, id); err != nil {
		log.Println(err.Error())
		return
	}
	templates.RenderTemplate(w, config.CategoryPostsTmpl, data)
}

