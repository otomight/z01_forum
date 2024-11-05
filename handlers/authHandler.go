package handlers

import (
	"Forum/database"
	"html/template"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

//// Registration \\\\

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		userName := r.FormValue("user_name")
		email := r.FormValue("email")
		password := r.FormValue("password")
		firstName := r.FormValue("first_name")
		lastName := r.FormValue("last_name")

		//Validate input
		if userName == "" || email == "" || password == "" || firstName == "" || lastName == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}

		// Hash password before saving it
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Password hashing error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Default user role
		userRole := "user"

		//Save user to database
		userID, err := database.SaveUser(userName, email, string(hashedPassword), firstName, lastName, userRole)
		if err != nil {
			log.Printf("Error saving user to database: %v", err)
			http.Error(w, "Unable to register user", http.StatusInternalServerError)
			return
		}

		//Log user automatically after registration
		sessionID, err := database.CreateUserSession(userID, userRole)
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
		return
	}
	renderRegistrationPage(w, r)
}

func renderRegistrationPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles("web/templates/register.html")
	if err != nil {
		http.Error(w, "Unable to render template:"+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

//// Login \\\\

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		//Validate User credentials
		user, err := database.ValidateUserCredentials(username, password)
		if err != nil {
			if err.Error() == "user not found" || err.Error() == "invalid password" {
				http.Error(w, "Invalid credentials", http.StatusUnauthorized)
				return
			} else {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		//Create new session for logged user
		sessionID, err := database.CreateUserSession(user.UserID, user.UserRole)
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
		return
	}

	// redirect to login page
	renderLogingPage(w, r)
}

func renderLogingPage(w http.ResponseWriter, r *http.Request) {
	//Get requests only
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Render Login template
	tmpl, err := template.ParseFiles("web/templates/login.html")
	if err != nil {
		http.Error(w, "Unable to render template:"+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
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
