package handlers

import (
	"fmt"
	"forum/internal/config"
	"forum/internal/server/services"
	"net/http"
)

// Google
func GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	authURL := fmt.Sprintf(
		"%s?client_id=%s&redirect_uri=%s&response_type=code&scope=email profile",
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
		ClientSecret: config.EnvVar.GoogleClientSecret,
		RedirectURI:  config.GoogleRedirectURI,
	}
	services.OAuthCallbackHandler(w, r, googleConfig)
}

// Discord
func DiscordLoginHandler(w http.ResponseWriter, r *http.Request) {
	authURL := fmt.Sprintf(
		"%s?client_id=%s&redirect_uri=%s&response_type=code&scope=identify email",
		config.DiscordAuthURL,
		config.DiscordClientID,
		config.DiscordRedirectURI,
	)
	http.Redirect(w, r, authURL, http.StatusFound)
}

func DiscordCallBackHandler(w http.ResponseWriter, r *http.Request) {
	discordConfig := config.ProviderConfig{
		Name:         "discord",
		TokenURL:     config.DiscordTokenURL,
		UserInfoURL:  config.DiscordUserInfoURL,
		ClientID:     config.DiscordClientID,
		ClientSecret: config.EnvVar.DiscordClientSecret,
		RedirectURI:  config.DiscordRedirectURI,
	}
	services.OAuthCallbackHandler(w, r, discordConfig)
}

// Fb
func FacebookLoginHandler(w http.ResponseWriter, r *http.Request) {
	authURL := fmt.Sprintf(
		"%s?client_id=%s&redirect_uri=%s&scope=email, public_profile",
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
		ClientSecret: config.EnvVar.FacebookClientSecret,
		RedirectURI:  config.FacebookRedirectURI,
	}
	services.OAuthCallbackHandler(w, r, facebookConfig)
}
