package handlers

import (
	"Forum/database"
	"fmt"
	"html/template"
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
	userID, _ := r.Context().Value(UserIDKey).(int)
	userRole, _ := r.Context().Value(UserRoleKey).(string)
	userName, _ := r.Context().Value(UserNameKey).(string)

	//Create data struct to hold posts
	data := struct {
		Title      string
		Posts      []database.Post
		IsLoggedIn bool
		UserID     int
		UserName   string
		UserRole   string
	}{
		Title:      fmt.Sprintf("Welcome, %s", userName),
		Posts:      posts,
		IsLoggedIn: userID != 0,
		UserID:     userID,
		UserName:   userName,
		UserRole:   userRole,
	}

	// Load all templates in the specified directory
	tmpl, err := template.ParseGlob("web/templates/*.html")
	if err != nil {
		http.Error(w, "Failed to load templates: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the base layout template with dynamic data
	if err := tmpl.ExecuteTemplate(w, "base_layout.html", data); err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// Registered Users Home Page depending on their roles
func RenderBaseHomePage(w http.ResponseWriter, r *http.Request) {
	userRole, _ := r.Context().Value(UserRoleKey).(string)

	//Choose template depending on role
	switch userRole {
	case "administrator":
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	case "moderator":
		http.Redirect(w, r, "/moderator", http.StatusSeeOther)
	case "user":
		http.Redirect(w, r, "/user", http.StatusSeeOther)
	default:
		RenderHomePage(w, r)
	}
}
