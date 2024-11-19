package handlers

import (
	"fmt"
	"forum/internal/config"
	"forum/internal/utils"
	"net/http"
)

// Google
func GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	authURL := fmt.Sprintf(
		"%s?client_id=%s&redirect_uri=%s&response_type=code&scope=openid profile email",
		config.GoogleAuthURL,
		config.GoogleClientID,
		config.GoogleRedirectURI,
	)
	http.Redirect(w, r, authURL, http.StatusFound)
}

func GoogleCallBackHandler(w http.ResponseWriter, r *http.Request) {
	googleConfig := config.ProviderConfig{
		Name:         "google",
		TokenURL:     config.GoogleTokenURL,
		UserInfoURL:  config.GoogleUserInfoURL,
		ClientID:     config.GoogleClientID,
		ClientSecret: config.GoogleClientSecret,
		RedirectURI:  config.GoogleRedirectURI,
	}
	utils.OAuthCallbackHandler(w, r, googleConfig)
}

// Github
func GithubLoginHandler(w http.ResponseWriter, r *http.Request) {
	authURL := fmt.Sprintf(
		"%s?client_id=%s&redirect_uri=%s&scope=user",
		config.GithubAuthURL,
		config.GithubClientID,
		config.GithubRedirectURI,
	)
	http.Redirect(w, r, authURL, http.StatusFound)
}

func GithubCallBackHandler(w http.ResponseWriter, r *http.Request) {
	githubConfig := config.ProviderConfig{
		Name:         "github",
		TokenURL:     config.GithubTokenURL,
		UserInfoURL:  config.GithubUserInfoURL,
		ClientID:     config.GithubClientID,
		ClientSecret: config.GithubClientSecret,
		RedirectURI:  config.GithubRedirectURI,
	}
	utils.OAuthCallbackHandler(w, r, githubConfig)
}

// Fb
func FacebookLoginHandler(w http.ResponseWriter, r *http.Request) {
	authURL := fmt.Sprintf(
		"%s?client_id=%s&redirect_uri=%s&scope=email",
		config.FacebookAuthURL,
		config.FacebookClientID,
		config.FacebookRedirectURI,
	)
	http.Redirect(w, r, authURL, http.StatusFound)
}

func FacebookCallBackHandler(w http.ResponseWriter, r *http.Request) {
	facebookConfig := config.ProviderConfig{
		Name:         "facebook",
		TokenURL:     config.FacebookTokenURL,
		UserInfoURL:  config.FacebookUserInfoURL,
		ClientID:     config.FacebookClientID,
		ClientSecret: config.FacebookClientSecret,
		RedirectURI:  config.FacebookRedirectURI,
	}
	utils.OAuthCallbackHandler(w, r, facebookConfig)
}
