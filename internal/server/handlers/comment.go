package handlers

import (
	"fmt"
	"forum/internal/config"
	"forum/internal/database"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Extract postID from the URL path
	urlPath := r.URL.Path
	parts := strings.Split(urlPath, "/")
	postIDStr := parts[len(parts)-1] // last part of the URL is the postID

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		log.Printf("Invalid postID: %v", err)
		http.Error(w, "Invalid postID", http.StatusBadRequest)
		return
	}

	session, ok := r.Context().Value(config.SessionKey).(*database.UserSession)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")
	if content == "" {
		http.Error(w, "Comment cannot be empty", http.StatusBadRequest)
		return
	}

	log.Printf("Attempting to add comment : userID=%d, postID=%d, content=%s",
												session.UserID, postID, content)

	err = database.AddComment(postID, session.UserID, content)
	if err != nil {
		log.Printf("Failed to add comment: %v", err)
		http.Error(w, "Failed to add comment", http.StatusInternalServerError)
		return
	}
	log.Println("Comment added successfully, redirecting to home page.")
	http.Redirect(w, r, fmt.Sprintf("/post/view/%d", postID), http.StatusSeeOther)
}
