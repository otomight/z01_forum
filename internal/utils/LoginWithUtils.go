package utils

import (
	"encoding/json"
	"fmt"
	"forum/internal/config"
	"forum/internal/database"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func ExtractCode(r *http.Request) (string, error) {
	code := r.URL.Query().Get("code")
	if code == "" {
		return "", fmt.Errorf("authorization code not found in request")
	}
	return code, nil
}

func ExchangeCodeForToken(tokenURL, clientID, clientSecret, redirectURI, code string) (string, error) {
	tokenReqBody := fmt.Sprintf(
		"code=%s&client_id=%s&client_secret=%s&redirect_uri=%s&grant_type=authorization_code",
		code, clientID, clientSecret, redirectURI,
	)

	resp, err := http.Post(tokenURL, "application/x-www-form-urlencoded", strings.NewReader(tokenReqBody))
	if err != nil {
		return "", fmt.Errorf("failed to exchange code for token: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read token response: %v", err)
	}

	var tokenResp map[string]interface{}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("failed to parse token response: %v", err)
	}

	accessToken, ok := tokenResp["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("access token not found in response")
	}
	return accessToken, nil
}

func FetchUserInfo(userInfoURL, accessToken string) (map[string]interface{}, error) {
	userResp, err := http.Get(fmt.Sprintf("%s?access_token=%s", userInfoURL, accessToken))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %v", err)
	}
	defer userResp.Body.Close()

	body, err := io.ReadAll(userResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read user info response: %v", err)
	}

	var userInfo map[string]interface{}
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %v", err)
	}
	return userInfo, nil
}

func getStringField(data map[string]interface{}, key string) (string, error) {
	// Check if the key exists in the map
	value, ok := data[key]
	if !ok {
		return "", fmt.Errorf("field '%s' is missing", key)
	}

	// Attempt to cast the value to a string
	strValue, ok := value.(string)
	if !ok || strValue == "" {
		return "", fmt.Errorf("field '%s' is either not a string or is empty", key)
	}

	return strValue, nil
}

func OAuthCallbackHandler(w http.ResponseWriter, r *http.Request, config config.ProviderConfig) {
	code, err := ExtractCode(r)
	if err != nil {
		log.Printf("Failed to extract code: %v", err)
		http.Error(w, "Failed to extract code:"+err.Error(), http.StatusBadRequest)
		return
	}

	accesToken, err := ExchangeCodeForToken(
		config.TokenURL,
		config.ClientID,
		config.ClientSecret,
		config.RedirectURI,
		code,
	)
	if err != nil {
		log.Printf("Failed to exchange code for token: %v", err)
		http.Error(w, "Failed to exchange code for token:"+err.Error(), http.StatusInternalServerError)
		return
	}

	userInfo, err := FetchUserInfo(config.UserInfoURL, accesToken)
	if err != nil {
		log.Printf("Failed to fetch user info: %v", err)
		http.Error(w, "Failed to fetch user info:"+err.Error(), http.StatusInternalServerError)
		return
	}

	// Safely extract fields from userInfo
	oauthID, err := getStringField(userInfo, "id")
	if err != nil {
		log.Printf("Invalid user info: %v", err)
		http.Error(w, "Invalid user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	email, _ := getStringField(userInfo, "email") // Email may be optional in some OAuth flows
	name, _ := getStringField(userInfo, "name")   // Use empty strings if missing
	avatar, _ := getStringField(userInfo, "picture")

	user, err := database.GetOrCreateUserByOAuth(
		config.Name,
		oauthID,
		email,
		name,
		avatar,
	)
	if err != nil {
		log.Printf("Failed to create/retrieve user: %v", err)
		http.Error(w, "Failed to create or retrieve user", http.StatusInternalServerError)
		return
	}

	sessionID, err := database.CreateUserSession(user.ID, user.UserRole, user.UserName)
	if err != nil {
		log.Printf("Failed to create session: %v", err)
		http.Error(w, "Failed to create session:"+err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
