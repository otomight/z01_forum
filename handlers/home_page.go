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

	//get userID and role from context
	userID, isLoggedIn := r.Context().Value(UserIDKey).(int)
	userRole, _ := r.Context().Value(UserRoleKey).(string)

	//Create data struct to hold posts
	data := struct {
		Title      string
		Posts      []database.Post
		IsLoggedIn bool
		UserID     int
		UserRole   string
	}{
		Title:      "Home",
		Posts:      posts,
		IsLoggedIn: isLoggedIn,
		UserID:     userID,
		UserRole:   userRole,
	}

	if err := tmpl.ExecuteTemplate(w, "home_page.html", data); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
