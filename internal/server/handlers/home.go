package handlers

import (
	"forum/internal/config"
	"forum/internal/database"
	"forum/internal/server/models"
	"forum/internal/server/templates"
	"log"
	"net/http"
	"time"
)

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve session ID from the cookie
	sessionCookie, err := r.Cookie("session_id")
	var session *database.UserSession
	var isLoggedIn bool
	var userName, userRole string
	var userID int

	// Check if session cookie is found
	if err != nil || sessionCookie.Value == "" {
		// No session cookie found: user is not logged in
		isLoggedIn = false
		userID = 0
		userName = ""
		userRole = ""
		log.Println("No session cookie found: user is not logged in")
	} else {
		// Session cookie exists: retrieve session from DB
		session, err = database.GetSessionByID(sessionCookie.Value)
		if err != nil || time.Now().After(session.Expiration) {
			// Session expired or invalid
			log.Println("Session expired or invalid")
			http.Error(w, "Session expired or invalid", http.StatusUnauthorized)
			return
		}

		// Session is valid, user is logged in
		isLoggedIn = true
		userID = session.UserID
		userRole = session.UserRole

		// Log the retrieved session info
		log.Printf("Session found for UserID: %d, UserRole: %s, UserName: %s", userID, userRole, userName)
	}

	// Prepare posts (assuming posts are public and do not depend on login)
	posts, err := database.GetAllPosts()
	if err != nil {
		log.Printf("Failed to retrieve posts: %v", err)
		http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError)
		return
	}

	// Prepare data for template rendering
	data := models.HomePageData{
		Title:      "Welcome to the Forum",
		Posts:      posts,
		IsLoggedIn: isLoggedIn,
		UserID:     userID,
		UserName:   userName,
		UserRole:   userRole,
	}
	templates.RenderTemplate(w, config.HomeTmpl, data)
}
