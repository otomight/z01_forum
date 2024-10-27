package handlers

import (
	"Forum/database"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// Hash password using bcrypt & save User to database
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var client database.Client
	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(client.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Could not hash password", http.StatusInternalServerError)
		return
	}
	client.Password = string(hashedPassword)

	if err := database.CreateClient(&client); err != nil {
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(client)
}

// Verify user's credentials and create session
func LoginUser(w http.ResponseWriter, r *http.Request) {
	//Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	//Check request method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//Decode request body into a struct
	var loginInfo struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginInfo); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	//Retrieve user by username or email
	user, err := database.GetClientByUsernameOrEmail(loginInfo.Login)
	if err != nil || user == nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	//Compare provided password with stored hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginInfo.Password)); err != nil {
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}

	//Create session for user
	sessionID, err := database.CreateUserSession(user.UserID)
	if err != nil {
		http.Error(w, "Could not create session", http.StatusInternalServerError)
		return
	}

	//Respond with sessionID + wanted details
	response := struct {
		SessionID string `json:"session_id"`
		UserID    int    `json:"user_id"`
		UserName  string `json:"user_name"`
	}{
		SessionID: sessionID,
		UserID:    user.UserID,
		UserName:  user.UserName,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
