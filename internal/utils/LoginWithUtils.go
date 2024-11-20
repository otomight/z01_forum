package utils

import (
	"encoding/json"
	"fmt"
	"forum/internal/config"
	"forum/internal/database"
	"io"
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

func OAuthCallbackHandler(w http.ResponseWriter, r *http.Request, config config.ProviderConfig) {
	code, err := ExtractCode(r)
	if err != nil {
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
		http.Error(w, "Failed to exchange code for token:"+err.Error(), http.StatusInternalServerError)
		return
	}

	userInfo, err := FetchUserInfo(config.UserInfoURL, accesToken)
	if err != nil {
		http.Error(w, "Failed to fetch user info:"+err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := database.GetOrCreateUserByOAuth(
		config.Name,
		userInfo["id"].(string),
		userInfo["email"].(string),
		userInfo["name"].(string),
		userInfo["picture"].(string),
	)
	if err != nil {
		http.Error(w, "Failed to create/Retrieve user:"+err.Error(), http.StatusInternalServerError)
		return
	}

	sessionID, err := database.CreateUserSession(user.UserID, user.UserRole, user.UserName)
	if err != nil {
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
