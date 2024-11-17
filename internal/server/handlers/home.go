package handlers

import (
	"forum/internal/config"
	"forum/internal/database"
	"forum/internal/server/models"
	"forum/internal/server/services"
	"forum/internal/server/templates"
	"net/http"
)

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
<<<<<<< HEAD
	var session	*database.UserSession
	var posts	[]database.Post
=======
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
		userName = session.UserName

		// Log the retrieved session info
		log.Printf("Session found for UserID: %d, UserRole: %s, UserName: %s", userID, userRole, userName)
	}
>>>>>>> 4e052d0 (added the userName to session table)

	session, _ = services.GetSession(r)
	posts, _ = database.GetAllPosts()
	// Prepare posts (assuming posts are public and do not depend on login)
	data := models.HomePageData{
		Posts:      posts,
		Session:	session,
	}
	templates.RenderTemplate(w, config.HomeTmpl, data)
}
