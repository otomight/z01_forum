package handlers

import (
	"Forum/database"
	"html/template"
	"net/http"
)

func AdminDashBoard(w http.ResponseWriter, r *http.Request) {
	userRole, ok := r.Context().Value(UserRoleKey).(string)
	if !ok || userRole != "administrator" {
		http.Error(w, "Acces denied", http.StatusForbidden)
		return
	}

	//Retrieve all posts
	posts, err := database.GetAllPosts()
	if err != nil {
		http.Error(w, "Failed to retrieve post", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title      string
		Posts      []database.Post
		IsLoggedIn bool
		UserRole   string
	}{
		Title:      "Admin dashboard",
		Posts:      posts,
		IsLoggedIn: true,
		UserRole:   userRole,
	}

	//Parse/execute admin dashboard template
	tmpl, err := template.ParseFiles("web/templates/logged_administrator_home_page.html")
	if err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Unable to render template", http.StatusInternalServerError)
	}
}
