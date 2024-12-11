package handlers

import (
	"forum/internal/config"
	"forum/internal/database"
	db "forum/internal/database"
	"forum/internal/server/models"
	"forum/internal/server/templates"
	"forum/internal/utils"
	"log"
	"net/http"
	"time"

	"github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

//// Registration \\\\

func createUser(
	w http.ResponseWriter, form models.RegisterForm, userRole string,
) (int, error) {
	var	hashedPassword	[]byte
	var	userID			int
	var	sqliteErr		sqlite3.Error
	var	tableName		string
	var	columnName		string
	var	cl				config.ClientsTableKeys
	var	ok				bool
	var	err				error

	cl = config.TableKeys.Clients
	hashedPassword, err = bcrypt.GenerateFromPassword(
		[]byte(form.Password), bcrypt.DefaultCost,
	)
	if err != nil {
		log.Printf("Password hashing error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return 0, err
	}
	userID, err = database.SaveUser(
		form.UserName, form.Email, string(hashedPassword), userRole,
	)
	if err == nil {
		return userID, nil
	}
	if sqliteErr, ok = err.(sqlite3.Error); ok {
		tableName, columnName = utils.GetSqlite3UniqueErrorInfos(sqliteErr)
		if tableName == cl.Clients && columnName == cl.UserName {
			http.Error(w, "Username already taken", http.StatusBadRequest)
			return 0, err
		} else if tableName == cl.Clients && columnName == cl.Email {
			http.Error(w, "Email already taken", http.StatusBadRequest)
			return 0, err
		}
	}
	log.Printf("Error saving user to database: %v", err)
	http.Error(w, "Unable to register user", http.StatusInternalServerError)
	return 0, err
}

func createSession(
	w http.ResponseWriter, userID int, userRole string, userName string,
) error {
	sessionID, err := database.CreateUserSession(userID, userRole, userName)
	if err != nil {
		log.Printf("Error creating session: %v", err)
		http.Error(w, "Failed to create sesion", http.StatusInternalServerError)
		return err
	}
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

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var	form		models.RegisterForm
	var	userID		int
	var	userRole	string
	var	session		*db.UserSession
	var	categories	[]*db.Category
	var	data		models.RegisterPageData
	var	err			error

	if r.Method != http.MethodPost {
		session, _ = r.Context().Value(config.SessionKey).(*db.UserSession)
		if session != nil {
			http.Error(w, "You are already logged.", http.StatusBadRequest)
			return
		}
		if categories, err = db.GetGlobalCategories(); err != nil {
			http.Error(
				w, "Error at fetching categories",
				http.StatusInternalServerError,
			)
		}
		data = models.RegisterPageData{
			Session:	session,
			Categories:	categories,
		}
		templates.RenderTemplate(w, config.RegisterTmpl, data)
		return
	}
	if err := utils.ParseForm(r, &form); err != nil {
		http.Error(w, "Unable to parse form:"+err.Error(),
			http.StatusBadRequest)
		return
	}
	if form.UserName == "" || form.Email == "" || form.Password == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}
	userRole = "user"
	if userID, err = createUser(w, form, userRole); err != nil {
		return
	}
	if err = createSession(w, userID, userRole, form.UserName); err != nil {
		return
	}
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
	var	form		models.LoginForm
	var	session		*database.UserSession
	var	categories	[]*db.Category
	var	data		models.LoginPageData
	var	err			error

	if r.Method != http.MethodPost {
		// redirect to login page
		session, _ = r.Context().Value(config.SessionKey).(*db.UserSession)
		if session != nil {
			http.Error(w, "You are already logged.", http.StatusBadRequest)
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
			Session:	session,
			Categories:	categories,
		}
		templates.RenderTemplate(w, config.LoginTmpl, data)
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
