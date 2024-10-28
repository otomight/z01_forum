package handlers

import (
	"html/template"
	"net/http"
)

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseGlob("./web/templates/*.html"))
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title          string
		IsLoginPage    bool
		IsRegisterPage bool
		IsLoggedIn     bool
	}{
		Title:          "Register",
		IsLoginPage:    false,
		IsRegisterPage: true,
		IsLoggedIn:     false, // Set based on user session
	}

	// Render the register.html template
	if err := tmpl.ExecuteTemplate(w, "register.html", data); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title          string
		IsLoginPage    bool
		IsRegisterPage bool
		IsLoggedIn     bool
	}{
		Title:          "Login",
		IsLoginPage:    true,
		IsRegisterPage: false,
		IsLoggedIn:     false, // Set based on user session
	}

	// Render the login.html template
	if err := tmpl.ExecuteTemplate(w, "login.html", data); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
