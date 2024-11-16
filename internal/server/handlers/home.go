package handlers

import (
	"fmt"
	"forum/internal/config"
	"forum/internal/database"
	"forum/internal/server/models"
	"forum/internal/server/services"
	"forum/internal/server/templates"
	"net/http"
)

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	var session	*database.UserSession
	var posts	[]database.Post

	session, err := services.GetSession(r)
	posts, _ = database.GetAllPosts()
	// Prepare posts (assuming posts are public and do not depend on login)
	fmt.Println("ERROR:", err)
	fmt.Println(session.IsLoggedIn)
	data := models.HomePageData{
		Posts:      posts,
		Session:	session,
	}
	templates.RenderTemplate(w, config.HomeTmpl, data)
}
