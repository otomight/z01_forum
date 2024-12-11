package handlers

import (
	"forum/internal/config"
	db "forum/internal/database"
	"forum/internal/server/models"
	"forum/internal/server/templates"
	"forum/internal/utils"
	"log"
	"net/http"
	"time"
)

func removeExistingUserSession(user db.Client) {
	var	session	*db.UserSession

	session, _ = db.GetSessionByUserID(user.ID)
	if session == nil {
		return // no session found or any other error
	}
	db.DeleteSession(session.ID)
}

func displayLoginPage(w http.ResponseWriter, r *http.Request) {
	var	session		*db.UserSession
	var	categories	[]*db.Category
	var	data		models.LoginPageData
	var	err			error

	session, _ = r.Context().Value(config.SessionKey).(*db.UserSession)
	if session != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if categories, err = db.GetGlobalCategories(); err != nil {
		http.Error(
			w, "Error at fetching categories",
			http.StatusInternalServerError,
		)
		return
	}
	data = models.LoginPageData{
		Session:	nil,
		Categories:	categories,
	}
	templates.RenderTemplate(w, config.LoginTmpl, data)
}

func createSessionOnLogged(w http.ResponseWriter, user *db.Client) error {
	var	sessionID	string
	var	err			error
	removeExistingUserSession(*user)
	sessionID, err = db.CreateUserSession(user.ID, user.UserRole, user.UserName)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return err
	}

	//Set sessionID in cookie
	http.SetCookie(w, &http.Cookie{
		Name:		"session_id",
		Value:		sessionID,
		Path:		"/",
		Expires:	time.Now().Add(24 * time.Hour),
		HttpOnly:	true,
		Secure:		true,
		SameSite:	http.SameSiteLaxMode,
	})
	return nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var	form		models.LoginForm
	var	err			error

	if r.Method == http.MethodGet {
		displayLoginPage(w, r)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if err := utils.ParseForm(r, &form); err != nil {
		http.Error(
			w, "Unable to parse form:"+err.Error(), http.StatusBadRequest,
		)
		return
	}
	log.Printf("Attempting to log in user: %s", form.Username)
	user, err := db.ValidateUserCredentials(form.Username, form.Password)
	if err != nil {
		if err.Error() == "user not found" || err.Error() == "invalid password" {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
	if err = createSessionOnLogged(w, &user); err != nil {
		return
	}
	log.Printf("User successfully logged in with role: %s", user.UserRole)
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func LogOutHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve session ID from cookie
	cookie, err := r.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "No session found", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Error retrieving session", http.StatusInternalServerError)
		return
	}

	err = db.DeleteSession(cookie.Value)
	if err != nil {
		http.Error(w, "Failed to log out", http.StatusInternalServerError)
		return
	}

	// Optionally, clear the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour), // Set an expiration in the past
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	// Redirect to unlogged home page after logout
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
