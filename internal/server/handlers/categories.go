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

func fillPostsInCategoriesPageData(
	r *http.Request,
	categoryID int,
) (*models.PostsInCategoriesPageData, error) {
	var	session		*db.UserSession
	var	data		*models.PostsInCategoriesPageData
	var	category	*db.Category
	var	posts		[]*db.Post
	var	err			error

	session, _ = r.Context().Value(config.SessionKey).(*db.UserSession)
	category, err = db.GetGlobalCategoryByID(categoryID)
	if err != nil {
		return nil, err
	}
	posts, err = db.GetPostsByCategoryID(categoryID)
	data = &models.PostsInCategoriesPageData{
		Session: session,
		Category: category,
		Posts: posts,
	}
	return data, nil
}

func PostsInCategoriesPageHandler(w http.ResponseWriter, r *http.Request) {
	var	idStr	string
	var	id		int
	var	err		error
	var	data	*models.PostsInCategoriesPageData

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
	data, err = fillPostsInCategoriesPageData(r, id)
	if err != nil {
		log.Println(err.Error())
		http.Error(
			w, "Failed to fetch posts from category",
			http.StatusInternalServerError,
		)
	}
	templates.RenderTemplate(w, config.ListPosts, data)
}

func CategoriesPageHandler(w http.ResponseWriter, r *http.Request) {
	var	session		*db.UserSession
	var	categories	[]*db.Category
	var	err			error

	if r.Method != http.MethodGet {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}
	session, _ = r.Context().Value(config.SessionKey).(*db.UserSession)
	categories, err = db.GetGlobalCategories()
	if err != nil {
		http.Error(w, "Error at fetching categories", http.StatusInternalServerError)
	}
	data := models.CategoriesPageData{
		Session:	session,
		Categories:	categories,
	}
	templates.RenderTemplate(w, config.CategoriesTmpl, data)
}

