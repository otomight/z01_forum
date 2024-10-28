package handlers

import (
	"Forum/database"
	"html/template"
	"net/http"
)

func ModeratorDashboardHAndler(w http.ResponseWriter, r *http.Request) {
	userRole, ok := r.Context().Value(UserRoleKey).(string)
	if !ok || userRole != "moderator" {
		http.Error(w, "Acces denied", http.StatusForbidden)
		return
	}

	posts, err := database.GetAllPosts()
	if err != nil {
		http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title      string
		Posts      []database.Post
		IsLoggedIn bool
		UserRole   string
	}{
		Title:      "Moderator dashboard",
		Posts:      posts,
		IsLoggedIn: true,
		UserRole:   userRole,
	}

	tmpl, err := template.ParseFiles("web/templates/logged_moderator_home_page.html")
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
	}
}
