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
	var client database.Client
	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	//Retrieve stored client data from the provided email address
	storedClient, err := database.GetClientByEmail(client.Email)
	if err != nil || storedClient == nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	//Compare provided password with stored hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(storedClient.Password), []byte(client.Password)); err != nil {
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}

	//Create session for user
	sessionID, err := database.CreateUserSession(storedClient.UserID)
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
		UserID:    storedClient.UserID,
		UserName:  storedClient.UserName,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
