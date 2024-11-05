package handlers

import (
	"Forum/database"
	"html/template"
	"log"
	"net/http"
)

var tmpl *template.Template

// Initialize template package + parse all HTML templates
func init() {
	tmpl = template.Must(template.ParseGlob("./web/templates/*.html"))
}

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	//Retrieve user infos from context
	userID, idOk := r.Context().Value(UserIDKey).(int)
	userRole, roleOk := r.Context().Value(UserRoleKey).(string)
	userName, nameOk := r.Context().Value(UserNameKey).(string)

	isLoggedIn := idOk && roleOk && nameOk && userID != 0

	posts, err := database.GetAllPosts()
	if err != nil {
		log.Printf("Failed to retrieve post: %v", err)
		http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError)
		return
	}

	//Prepare data for template rendering
	data := struct {
		Title      string
		Posts      []database.Post
		IsLoggedIn bool
		UserID     int
		UserName   string
		UserRole   string
	}{
		Title:      "Welcome to the Forum",
		Posts:      posts,
		IsLoggedIn: isLoggedIn,
		UserID:     userID,
		UserName:   userName,
		UserRole:   userRole,
	}

	// Execute the home_page template with dynamic data
	if err := tmpl.ExecuteTemplate(w, "home_page.html", data); err != nil {
		log.Printf("Template execution error in 'home_page': %v", err)
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
