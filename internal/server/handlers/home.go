package handlers

import (
	"forum/internal/config"
	db "forum/internal/database"
	"forum/internal/server/models"
	"forum/internal/server/templates"
	"log"
	"net/http"
)

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	var	session	*db.UserSession
	var	posts	[]*db.Post
	var	userID	int
	var	err		error

	session, _ = r.Context().Value(config.SessionKey).(*db.UserSession)
	if session == nil {
		userID = 0
	} else {
		userID = session.UserID
	}
	posts, err = db.GetAllPosts(userID)
	if err != nil {
		log.Printf("Failed to retrieve posts: %v", err)
		http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError)
		return
	}
	data := models.HomePageData{
		Session:	session,
		Posts:		posts,
	}
	templates.RenderTemplate(w, config.HomeTmpl, data)
}
