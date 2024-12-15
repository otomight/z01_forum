package handlers

import (
	"forum/internal/config"
	"forum/internal/database"
	db "forum/internal/database"
	"forum/internal/server/models"
	"forum/internal/server/templates"
	sessioncreate "forum/internal/sessioncreate"
	"forum/internal/utils"
	"log"
	"net/http"

	"github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func displayRegisterPage(
	w http.ResponseWriter, r *http.Request,
	userInput *models.RegisterPageUserInput, errorMsg *models.RegisterErrorMsg,
) {
	var session *db.UserSession
	var categories []*db.Category
	var data models.RegisterPageData
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
	data = models.RegisterPageData{
		Session:    nil,
		Categories: categories,
		UserInput:  userInput,
		ErrorMsg:   errorMsg,
	}
	templates.RenderTemplate(w, config.RegisterTmpl, data)
}

func handlerRegisterError(
	w http.ResponseWriter, r *http.Request,
	err error, form models.RegisterForm,
) {
	var sqliteErr sqlite3.Error
	var tableName string
	var columnName string
	var cl config.ClientsTableKeys
	var userInput *models.RegisterPageUserInput
	var ok bool

	cl = config.TableKeys.Clients
	if sqliteErr, ok = err.(sqlite3.Error); ok {
		tableName, columnName = utils.GetSqlite3UniqueErrorInfos(sqliteErr)
		userInput = &models.RegisterPageUserInput{
			Username: form.UserName,
			Email:    form.Email,
		}
		if tableName == cl.Clients && columnName == cl.UserName {
			displayRegisterPage(w, r, userInput, &models.RegisterErrorMsg{
				UsernameAlreadyTaken: "User name already taken",
			})
			return
		} else if tableName == cl.Clients && columnName == cl.Email {
			displayRegisterPage(w, r, userInput, &models.RegisterErrorMsg{
				EmailAlreadyTaken: "Email already taken",
			})
			return
		}
	}
	log.Printf("Error saving user to database: %v", err)
	http.Error(w, "Unable to register user", http.StatusInternalServerError)
}

func createUser(
	w http.ResponseWriter, r *http.Request,
	form models.RegisterForm, userRole string,
) (int, error) {
	var hashedPassword []byte
	var userID int
	var err error

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
	handlerRegisterError(w, r, err, form)
	return 0, err
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var form models.RegisterForm
	var userID int
	var err error

	if r.Method == http.MethodGet {
		displayRegisterPage(w, r, nil, nil)
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
	if userID, err = createUser(w, r, form, config.UserRole.User); err != nil {
		return
	}
	err = sessioncreate.SessionCreate(w, userID, config.UserRole.User, form.UserName)
	if err != nil {
		return
	}
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
