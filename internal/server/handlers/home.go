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
	var session	*db.UserSession
	var posts	[]db.Post
	var	err		error

	session, _ = r.Context().Value(config.SessionKey).(*db.UserSession)
	posts, err = db.GetAllPosts()
	// Prepare posts
	if err != nil {
		log.Printf("Failed to retrieve posts: %v", err)
		http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError)
		return
	}

	//For each post display the corresponding comments
	for i := range posts {
		comments, err := db.GetCommentsByPostID(posts[i].PostID)
		if err != nil {
			log.Printf("failed to fetch comments for post %d: %v", posts[i].PostID, err)
			continue
		}
		posts[i].Comments = comments
	}
	data := models.HomePageData{
		Posts:		posts,
		Session:	session,
	}
	templates.RenderTemplate(w, config.HomeTmpl, data)
}
