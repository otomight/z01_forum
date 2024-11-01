package handlers

import (
	"Forum/database"
	"html/template"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseGlob("./web/templates/*.html"))
}

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
		sessionID, err := database.CreateUserSession(userID, "user")
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

		// Redirect to base_layout.html if user is not an admin or moderator
		if userRole != "administrator" && userRole != "moderator" {
			http.Redirect(w, r, "/base_layout", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
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

	if err := tmpl.ExecuteTemplate(w, "register.html", nil); err != nil {
		http.Error(w, "Unable to render template: "+err.Error(), http.StatusInternalServerError)
		return
	}
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

		// Redirect based on user role
		switch user.UserRole {
		case "administrator":
			// Redirect to admin homepage
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
		case "moderator":
			// Redirect to moderator homepage
			http.Redirect(w, r, "/moderator", http.StatusSeeOther)
		default:
			// Redirect to base layout for regular users
			http.Redirect(w, r, "/base_layout", http.StatusSeeOther)
		}
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

	// Use the global tmpl variable to execute the login template
	if err := tmpl.ExecuteTemplate(w, "login.html", nil); err != nil {
		http.Error(w, "Unable to render template: "+err.Error(), http.StatusInternalServerError)
	}
}
