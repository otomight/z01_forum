package services

import (
	"encoding/json"
	"fmt"
	"forum/internal/config"
	"forum/internal/database"
	"forum/internal/sessioncreate"

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

func ExchangeCodeForToken(tokenURL, clientID, clientSecret, redirectURI, code, method string) (string, error) {
	// Default to "POST" if no method is provided
	if method == "" {
		method = "POST"
	}

	var resp *http.Response
	var err error

	// Check method and perform the correct HTTP request
	if method == "POST" {
		// Prepare the request body for POST
		tokenReqBody := fmt.Sprintf(
			"code=%s&client_id=%s&client_secret=%s&redirect_uri=%s&grant_type=authorization_code",
			code, clientID, clientSecret, redirectURI,
		)
		resp, err = http.Post(tokenURL, "application/x-www-form-urlencoded", strings.NewReader(tokenReqBody))
	} else if method == "GET" {
		// Prepare GET parameters
		reqURL := fmt.Sprintf("%s?code=%s&client_id=%s&client_secret=%s&redirect_uri=%s&grant_type=authorization_code",
			tokenURL, code, clientID, clientSecret, redirectURI)
		resp, err = http.Get(reqURL)
	} else {
		return "", fmt.Errorf("unsupported method: %s", method)
	}

	if err != nil {
		return "", fmt.Errorf("failed to exchange code for token: %v", err)
	}
	defer resp.Body.Close()

	var tokenResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("error decoding body : %w", err)
	}

	response, _ := json.MarshalIndent(tokenResp, "", " ")
	log.Printf("data recu : %s", response)

	accessToken, ok := tokenResp["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("access token not found in response")
	}
	return accessToken, nil
}

func FetchUserInfo(userInfoURL, accessToken, method string) (map[string]interface{}, error) {
	var userResp *http.Response
	var err error

	// Default to "GET" if no method is provided
	if method == "" {
		method = "GET"
	}

	// Prepare the request based on method
	if method == "GET" {
		// Use GET method to fetch user info (facebook)
		url := fmt.Sprintf("%s?access_token=%s", userInfoURL, accessToken)
		userResp, err = http.Get(url)
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
		// If there was an error during the request, return it
		return nil, fmt.Errorf("failed to fetch user info: %v", err)
	}
	defer userResp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(userResp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("error decoding body: %v", err) // Handle the error if decoding fails
	}

	response, _ := json.MarshalIndent(userInfo, "", " ")
	log.Printf("info user recu : %s", response)

	// Extract user ID
	var userID string
	if userInfo["id"] != nil {
		userID = userInfo["id"].(string) // Facebook
	} else if userInfo["sub"] != nil {
		userID = userInfo["sub"].(string) // Google
	} else {
		return nil, fmt.Errorf("missing required user ID field")
	}
	log.Printf("extracted user ID: %s", userID)

	return userInfo, nil
}

func OAuthCallbackHandler(w http.ResponseWriter, r *http.Request, config config.ProviderConfig) {
	code, err := ExtractCode(r)
	if err != nil {
		log.Printf("Failed to extract code: %v", err)
		http.Error(w, "Failed to extract code: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Determine method based on OAuth provider
	var method string
	switch config.Name {
	case "facebook":
		method = "GET" // Facebook
	case "google":
		method = "POST" // Google
	default:
		method = "POST" // Default to POST if the provider is not recognized
	}

	accesToken, err := ExchangeCodeForToken(
		config.TokenURL,
		config.ClientID,
		config.ClientSecret,
		config.RedirectURI,
		code,
		method,
	)
	if err != nil {
		log.Printf("Failed to exchange code for token: %v", err)
		http.Error(w, "Failed to exchange code for token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	userInfo, err := FetchUserInfo(config.UserInfoURL, accesToken, method)
	if err != nil {
		log.Printf("Failed to fetch user info: %v", err)
		http.Error(w, "Failed to fetch user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract user id and email, falling back to empty strings if not available
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

	// Map required fields to the Client struct
	client := &database.Client{
		OauthProvider: config.Name,
		OauthID:       id,
		Email:         email,
		UserName:      userInfo["name"].(string),
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
	if err = sessioncreate.SessionCreate(w, client.ID, client.UserRole, client.UserName); err != nil {
		log.Printf("Error creating session: %v", err)
		http.Error(w, "Failed to create session.", http.StatusInternalServerError)
		return
	}

	// Log and respond
	response, _ := json.MarshalIndent(dbClient, "", " ")
	log.Printf("info user recu: %s", response)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
