package handlers

import (
	"html/template"
	"net/http"
)

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseGlob("template/*.html"))
}

func LogRegisterHandler(w http.ResponseWriter, r *http.Request) {
	isRegisterPage := r.URL.Query().Get("register") == "true"

	data := struct {
		Title          string
		Content        string
		IsLoginPage    bool
		IsRegisterPage bool
		IsLoggedIn     bool
	}{
		Title:          "Login/Register",
		Content:        "",
		IsLoginPage:    !isRegisterPage,
		IsRegisterPage: isRegisterPage,
		IsLoggedIn:     false, // Set based on user session
	}

	// Render the log_register.html template with the data
	if err := tmpl.ExecuteTemplate(w, "log_register.html", data); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
