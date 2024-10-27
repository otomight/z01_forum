package handlers

import (
	"Forum/database"
	"net/http"
)

// Unauthenticated users home page
func RenderHomePage(w http.ResponseWriter, r *http.Request) {
	posts, err := database.GetAllPosts()
	if err != nil {
		http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError)
		return
	}

	//Create data struct to hold posts
	data := struct {
		Title      string
		Posts      []database.Post
		IsLoggedIn bool
		UserRole   string
	}{
		Title:      "Home",
		Posts:      posts,
		IsLoggedIn: false,
		UserRole:   "",
	}

	if err := tmpl.ExecuteTemplate(w, "home_page.html", data); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
