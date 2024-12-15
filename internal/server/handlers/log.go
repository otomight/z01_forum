package handlers

import (
	"forum/internal/config"
	db "forum/internal/database"
	"forum/internal/server/models"
	"forum/internal/server/templates"
	sessioncreate "forum/internal/sessioncreate"
	"forum/internal/utils"
	"log"
	"net/http"
	"time"
)

func removeExistingUserSession(user db.Client) {
	var session *db.UserSession

	session, _ = db.GetSessionByUserID(user.ID)
	if session == nil {
		return // no session found or any other error
	}
	db.DeleteSession(session.ID)
}

func displayLoginPage(
	w http.ResponseWriter, r *http.Request,
	userInput *models.LoginPageUserInput, errorMsg *models.LoginErrorMsg,
) {
	var session *db.UserSession
	var categories []*db.Category
	var data models.LoginPageData
	var err error

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
		Session:    nil,
		Categories: categories,
		UserInput:  userInput,
		ErrorMsg:   errorMsg,
	}
	templates.RenderTemplate(w, config.LoginTmpl, data)
}

func handleLogError(w http.ResponseWriter, r *http.Request, err error, username string) {
	var userInput *models.LoginPageUserInput

	userInput = &models.LoginPageUserInput{
		Username: username,
	}
	if err.Error() == "user not found" {
		displayLoginPage(w, r, userInput, &models.LoginErrorMsg{
			UserNotFound: "User not found",
		})
	} else if err.Error() == "invalid password" {
		displayLoginPage(w, r, userInput, &models.LoginErrorMsg{
			IncorrectPassword: "Invalid password",
		})
	} else {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func logUser(
	w http.ResponseWriter, r *http.Request,
	username string, password string,
) error {
	var user db.Client
	var err error

	user, err = db.ValidateUserCredentials(username, password)
	if err != nil {
		handleLogError(w, r, err, username)
		return err
	}
	removeExistingUserSession(user)
	err = sessioncreate.SessionCreate(w, user.ID, user.UserRole, user.UserName)
	if err != nil {
		return err
	}
	log.Printf("User successfully logged in with role: %s", user.UserRole)
	return nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var form models.LoginForm
	var err error

	if r.Method == http.MethodGet {
		displayLoginPage(w, r, nil, nil)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if err = utils.ParseForm(r, &form); err != nil {
		http.Error(
			w, "Unable to parse form:"+err.Error(), http.StatusBadRequest,
		)
		return
	}
	log.Printf("Attempting to log in user: %s", form.Username)
	if err = logUser(w, r, form.Username, form.Password); err != nil {
		return
	}
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func LogOutHandler(w http.ResponseWriter, r *http.Request) {
	var cookie *http.Cookie
	var err error

	if cookie, err = r.Cookie("session_id"); err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "No session found", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Error retrieving session", http.StatusInternalServerError)
		return
	}
	if err = db.DeleteSession(cookie.Value); err != nil {
		http.Error(w, "Failed to log out", http.StatusInternalServerError)
		return
	}
	// delete the cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session_id",
		Value:   "",
		Path:    "/",
		Expires: time.Now().Add(-1 * time.Hour),
		// Set an expiration in the past to delete the cookie
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
