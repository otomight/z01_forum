package handlers

import (
	"Forum/database"
	"html/template"
	"log"
	"net/http"
)

func UserHomePageHandler(w http.ResponseWriter, r *http.Request) {
	userRole, roleOk := r.Context().Value(UserRoleKey).(string)
	userID, idOk := r.Context().Value(UserIDKey).(int)

	if !roleOk || !idOk {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	renderTemplate(w, "home_page.html", userRole, userID)
}

func ModeratorHomePageHandler(w http.ResponseWriter, r *http.Request) {
	userRole, roleOk := r.Context().Value(UserRoleKey).(string)
	userID, idOk := r.Context().Value(UserIDKey).(int)

	if !roleOk || !idOk || userRole != "moderator" {
		http.Error(w, "Acces denied", http.StatusForbidden)
		return
	}

	renderTemplate(w, "home_page.html", userRole, userID)
}

func AdministratorHomePageHandler(w http.ResponseWriter, r *http.Request) {
	userRole, roleOk := r.Context().Value(UserRoleKey).(string)
	userID, idOk := r.Context().Value(UserIDKey).(int)

	if !roleOk || !idOk || userRole != "administrator" {
		http.Error(w, "Acces denied", http.StatusForbidden)
		return
	}

	renderTemplate(w, "base_layout.html", userRole, userID)
}

func renderTemplate(w http.ResponseWriter, tmpl string, userRole string, userID int) {
	// Fetch posts from the database
	posts, err := database.GetAllPosts()
	if err != nil {
		http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError)
		return
	}

	// Prepare data to pass to the templates
	data := struct {
		Title      string
		Posts      []database.Post
		IsLoggedIn bool
		UserRole   string
		UserID     int
	}{
		Title:      "Home Page",
		Posts:      posts,
		IsLoggedIn: true,
		UserRole:   userRole,
		UserID:     userID,
	}

	// Define paths to base and specific templates
	tmplPath := "web/templates/" + tmpl
	baseTmplPath := "web/templates/base_layout.html"

	// Parse templates with base layout
	templates, err := template.ParseFiles(baseTmplPath, tmplPath)
	if err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
		log.Printf("Template parsing error: %v", err)
		return
	}

	// Render template
	err = templates.ExecuteTemplate(w, "base_layout", data)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		log.Printf("Template execution error on 'base_layout': %v", err)
	}
}
