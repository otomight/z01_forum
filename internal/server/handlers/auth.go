package handlers

import (
	"forum/internal/config"
	"forum/internal/database"
	"forum/internal/server/models"
	"forum/internal/server/templates"
	"forum/internal/utils"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

//// Registration \\\\

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var form models.RegisterForm

	if r.Method != http.MethodPost {
		templates.RenderTemplate(w, config.RegisterTmpl, nil)
		return
	}
	// store form
	if err := utils.ParseForm(r, &form); err != nil {
		http.Error(w, "Unable to parse form:"+err.Error(),
			http.StatusBadRequest)
		return
	}

	//Validate input
	if form.UserName == "" || form.Email == "" || form.Password == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Hash password before saving it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password),
		bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Password hashing error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Default user role
	userRole := "user"

	//Save user to database
	userID, err := database.SaveUser(form.UserName, form.Email,
		string(hashedPassword), userRole)
	if err != nil {
		log.Printf("Error saving user to database: %v", err)
		http.Error(w, "Unable to register user", http.StatusInternalServerError)
		return
	}

	//Log user automatically after registration
	sessionID, err := database.CreateUserSession(userID, userRole, form.UserName)
	if err != nil {
		log.Printf("Error creating session: %v", err)
		http.Error(w, "Failed to create sesion", http.StatusInternalServerError)
		return
	}

	//Set sessionID in cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	// Redirect to /home after registration
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

//// Login \\\\

func removeExistingUserSession(user database.Client) {
	var	session	*database.UserSession

	session, _ = database.GetSessionByUserID(user.ID)
	if session == nil {
		return // no session found or any other error
	}
	database.DeleteSession(session.ID)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var form models.LoginForm
	if r.Method != http.MethodPost {
		// redirect to login page
		templates.RenderTemplate(w, config.LoginTmpl, nil)
		return
	}
	// store form
	if err := utils.ParseForm(r, &form); err != nil {
		http.Error(w, "Unable to parse form:"+err.Error(),
			http.StatusBadRequest)
		return
	}
	log.Printf("Attempting to log in user: %s", form.Username)

	//Validate User credentials
	user, err := database.ValidateUserCredentials(form.Username, form.Password)
	if err != nil {
		if err.Error() == "user not found" || err.Error() == "invalid password" {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	log.Printf("User successfully logged in with role: %s", user.UserRole)

	//Create new session for logged user
	removeExistingUserSession(user)
	sessionID, err := database.CreateUserSession(user.ID, user.UserRole, user.UserName)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	//Set sessionID in cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	// Redirect to /home
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

	err = database.DeleteSession(cookie.Value)
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
