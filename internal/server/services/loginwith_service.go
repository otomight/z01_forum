package services

import (
	"encoding/json"
	"fmt"
	"forum/internal/config"
	"forum/internal/database"

	"log"
	"net/http"
	"strings"
)

func ExtractCode(r *http.Request) (string, error) {
	code := r.URL.Query().Get("code")
	if code == "" {
		return "", fmt.Errorf("authorization code not found in request")
	}
	return code, nil
}

func exchangeCodeForToken(tokenURL, clientID, clientSecret, redirectURI, code string) (string, error) {
	// Prepare the request body for POST
	tokenReqBody := fmt.Sprintf(
		"code=%s&client_id=%s&client_secret=%s&redirect_uri=%s&grant_type=authorization_code",
		code, clientID, clientSecret, redirectURI,
	)

	// Use POST method to exchange the code for a token
	resp, err := http.Post(tokenURL, "application/x-www-form-urlencoded", strings.NewReader(tokenReqBody))
	if err != nil {
		return "", fmt.Errorf("failed to exchange code for token: %v", err)
	}
	defer resp.Body.Close()

	// Decode the response body to extract the token
	var tokenResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("error decoding body: %w", err)
	}

	response, _ := json.MarshalIndent(tokenResp, "", " ")
	log.Printf("data recu : %s", response)

	// Extract the access token
	accessToken, ok := tokenResp["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("access token not found in response")
	}

	return accessToken, nil
}

func fetchUserInfo(userInfoURL, accessToken, method string) (map[string]interface{}, error) {
	var userResp *http.Response
	var err error

	// Default to "GET" if no method is provided
	if method == "" {
		method = "GET"
	}

	if method == "GET" {
		// Use GET method to fetch user info (Facebook/Discord)
		req, err := http.NewRequest("GET", userInfoURL, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %v", err)
		}
		// Set Authorization header with Bearer token
		req.Header.Set("Authorization", "Bearer "+accessToken)
		req.Header.Del("Content-Type")

		// Execute the request
		client := &http.Client{}
		userResp, err = client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch user info: %v", err)
		}
	} else if method == "POST" {
		// Use POST method to fetch user info (google)
		userResp, err = http.Post(
			userInfoURL,
			"application/x-www-form-urlencoded",
			strings.NewReader(fmt.Sprintf("access_token=%s", accessToken)),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch user info: %v", err)
		}
	} else {
		return nil, fmt.Errorf("unsupported method: %s", method)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %v", err)
	}
	defer userResp.Body.Close()

	log.Printf("Response status code: %d", userResp.StatusCode)

	// // If the status code isn't 200 (OK), log and return an error
	// if userResp.StatusCode != 200 {
	// 	body, _ := io.ReadAll(userResp.Body)
	// 	log.Printf("Response body on error: %s", string(body))
	// 	return nil, fmt.Errorf("unexpected status code: %d", userResp.StatusCode)
	// }

	// Parse the JSON response body into a map
	var userInfo map[string]interface{}
	if err := json.NewDecoder(userResp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("error decoding body: %v", err)
	}

	// Log the received user information for debugging
	response, _ := json.MarshalIndent(userInfo, "", " ")
	log.Printf("Received user info: %s", response)

	// Extract user ID
	var userID string
	if userInfo["id"] != nil {
		userID = userInfo["id"].(string) // Facebook/Discord
	} else if userInfo["sub"] != nil {
		userID = userInfo["sub"].(string) // Google
	} else {
		return nil, fmt.Errorf("missing required user ID field")
	}
	log.Printf("Extracted user ID: %s", userID)

	return userInfo, nil
}

func OAuthCallbackHandler(w http.ResponseWriter, r *http.Request, config config.ProviderConfig) {
	code, err := ExtractCode(r)
	if err != nil {
		log.Printf("Failed to extract code: %v", err)
		http.Error(w, "Failed to extract code: "+err.Error(), http.StatusBadRequest)
		return
	}

	// log.Printf("Received authorization code: %s", code)

	accesToken, err := exchangeCodeForToken(
		config.TokenURL,
		config.ClientID,
		config.ClientSecret,
		config.RedirectURI,
		code,
	)
	if err != nil {
		log.Printf("Failed to exchange code for token: %v", err)
		http.Error(w, "Failed to exchange code for token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Method for fetching user info
	var userInfoMethod string
	switch config.Name {
	case "google":
		userInfoMethod = "POST" // Google t
	case "facebook", "discord":
		userInfoMethod = "GET" // Facebook/Discord
	default:
		userInfoMethod = "GET"
	}

	userInfo, err := fetchUserInfo(config.UserInfoURL, accesToken, userInfoMethod)
	if err != nil {
		log.Printf("Failed to fetch user info: %v", err)
		http.Error(w, "Failed to fetch user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, idExists := userInfo["id"].(string)
	if !idExists || id == "" {
		log.Printf("Id not provided by user, continuing without it")
		id = ""
	}

	email, emailExists := userInfo["email"].(string)
	if !emailExists || email == "" {
		log.Printf("Email not provided by user, will continue without it.")
		email = ""
	}

	username, usernameExists := userInfo["username"].(string)
	if !usernameExists || username == "" {
		log.Printf("Username not found in response, defaulting to empty string")
		username = "" // Or handle as appropriate
	}

	// Map required fields to the Client struct
	client := &database.Client{
		OauthProvider: config.Name,
		OauthID:       id,
		Email:         email,
		UserName:      username,
	}

	// Retrieve or create the client in the database
	dbClient, err := database.GetOrCreateUserByOAuth(
		client.OauthProvider,
		client.OauthID,
		client.Email,
		client.UserName,
		client.Avatar,
	)
	if err != nil {
		log.Printf("Failed to create/retrieve client: %v", err)
		http.Error(w, "Failed to create or retrieve user.", http.StatusInternalServerError)
		return
	}

	// Populate additional client fields
	client.ID = dbClient.ID

	// Create session
	if err = SessionCreate(w, client.ID, client.UserRole, client.UserName); err != nil {
		log.Printf("Error creating session: %v", err)
		http.Error(w, "Failed to create session.", http.StatusInternalServerError)
		return
	}

	// Log and respond
	response, _ := json.MarshalIndent(dbClient, "", " ")
	log.Printf("info user recu: %s", response)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
