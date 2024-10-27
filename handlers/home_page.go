package handlers

import (
	"Forum/database"
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

	tmpl := template.Must(template.ParseFiles("web/templates/home_page.html"))

	//Create data struct to hold posts
	data := struct {
		Posts []database.Post
	}{
		Posts: posts,
	}

	//Render template + pass in posts
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}
