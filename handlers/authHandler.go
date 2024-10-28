package handlers

import (
	"Forum/database"
	"html/template"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseGlob("./web/templates/*.html"))
}

func RegisterHAndler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		userName := r.FormValue("username")
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
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		//Default user role
		userRole := "user"

		//Save user to database
		userID, err := database.SaveUser(userName, email, string(hashedPassword), firstName, lastName, userRole)
		if err != nil {
			http.Error(w, "Unable to register user", http.StatusInternalServerError)
			return
		}

		//Log user automatically after registration
		sessionID, err := database.CreateUserSession(userID, userRole)
		if err != nil {
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

		//Redirect to logged_user_homepage after success
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Render the register.html template
	renderRegistrationPage(w, r)
}

func renderRegistrationPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles("web/templates/register.html")
	if err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Unable to render template", http.StatusInternalServerError)
		return
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		//Validate User credentials
		user, err := database.ValidateUserCredentials(username, password)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if user.UserID == 0 {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
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

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Show login.html template
	renderLogingPage(w, r)
}

func renderLogingPage(w http.ResponseWriter, r *http.Request) {
	//Get requests only
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//Parse template
	tmpl, err := template.ParseFiles("web/templates/login.html")
	if err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
		return
	}

	//Execute template
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Unable to render template", http.StatusInternalServerError)
	}
}
