package handlers

import (
	"forum/internal/config"
	db "forum/internal/database"
	"forum/internal/server/models"
	"forum/internal/server/templates"
	"net/http"
)

func CategoriesPageHandler(w http.ResponseWriter, r *http.Request) {
	var	session		*db.UserSession
	var	categories	[]*db.Category
	var	err			error

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

