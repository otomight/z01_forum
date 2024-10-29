package handlers

import (
	"Forum/database"
	"html/template"
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
	userID, idOk := r.Context().Value(UserRoleKey).(int)

	if !roleOk || !idOk || userRole != "moderator" {
		http.Error(w, "Acces denied", http.StatusForbidden)
		return
	}

	renderTemplate(w, "home_page.html", userRole, userID)
}

func AdministratorHomePageHandler(w http.ResponseWriter, r *http.Request) {
	userRole, roleOk := r.Context().Value(UserRoleKey).(string)
	userID, idOk := r.Context().Value(UserRoleKey).(int)

	if !roleOk || !idOk || userRole != "administrator" {
		http.Error(w, "Acces denied", http.StatusForbidden)
		return
	}

	renderTemplate(w, "home_page.html", userRole, userID)
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

	// Parse the base layout and the specific template together
	tmplPath := "web/templates/" + tmpl
	baseTmplPath := "web/templates/base_layout.html"
	templates := template.Must(template.ParseFiles(baseTmplPath, tmplPath))

	// Render the base layout with the dynamic data
	if err := templates.ExecuteTemplate(w, "base_layout", data); err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
	}
}
