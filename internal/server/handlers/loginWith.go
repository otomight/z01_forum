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
	code, err := utils.ExtractCode(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	accesToken, err := utils.ExchangeCodeForToken(
		config.GoogleTokenURL,
		config.GoogleClientID,
		config.GoogleClientSecret,
		config.GoogleRedirectURI,
		code,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userInfo, err := utils.FetchUserInfo(config.GoogleUserInfoURL, accesToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Hello %s, your email is %s", userInfo["name"], userInfo["email"])
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
	code, err := utils.ExtractCode(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	accesToken, err := utils.ExchangeCodeForToken(
		config.GithubTokenURL,
		config.GithubClientID,
		config.GithubClientSecret,
		config.GithubRedirectURI,
		code,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userInfo, err := utils.FetchUserInfo(config.GithubUserInfoURL, accesToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Hello %s your email is %s", userInfo["name"], userInfo["email"])
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
	code, err := utils.ExtractCode(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	accesToken, err := utils.ExchangeCodeForToken(
		config.FacebookTokenURL,
		config.FacebookClientID,
		config.FacebookClientSecret,
		config.FacebookRedirectURI,
		code,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userInfo, err := utils.FetchUserInfo(config.FacebookUserInfoURL, accesToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Hello %v, your email is %s", userInfo["name"], userInfo["email"])
}
