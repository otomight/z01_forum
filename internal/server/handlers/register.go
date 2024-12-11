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

func displayRegisterPage(w http.ResponseWriter, r *http.Request) {
	var	session		*db.UserSession
	var	categories	[]*db.Category
	var	data		models.RegisterPageData
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
	data = models.RegisterPageData{
		Session:	nil,
		Categories:	categories,
	}
	templates.RenderTemplate(w, config.RegisterTmpl, data)
}

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
	var	err			error

	if r.Method == http.MethodGet {
		displayRegisterPage(w, r)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Bad request", http.StatusBadRequest)
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
	if userID, err = createUser(w, form, config.UserRole.User); err != nil {
		return
	}
	err = createSession(w, userID, config.UserRole.User, form.UserName)
	if err != nil {
		return
	}
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
