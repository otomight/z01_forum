package handlers

import (
	"fmt"
	"forum/internal/config"
	"net/http"
)

func GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	authURL := fmt.Sprintf(
		"%s?client_id=%s&redirect_uri=%s&response_type=code&scope=openid profile email",
		config.GoogleAuthURL,
		config.GoogleClientID,
		config.GoogleRedirectURI,
	)
	http.Redirect(w, r, authURL, http.StatusFound)
}

func GithubLoginHandler(w http.ResponseWriter, r *http.Request) {
	authURL := fmt.Sprintf(
		"%s?client_id=%s&redirect_uri=%s&scope=user",
		config.GithubAuthURL,
		config.GithubClientID,
		config.GithubRedirectURI,
	)
	http.Redirect(w, r, authURL, http.StatusFound)
}

func FacebookLoginHandler(w http.ResponseWriter, r *http.Request) {
	authURL := fmt.Sprintf(
		"%s?client_id=%s&redirect_uri=%s&scope=email",
		config.FacebookAuthURL,
		config.FacebookClientID,
		config.FacebookRedirectURI,
	)
	http.Redirect(w, r, authURL, http.StatusFound)
}
