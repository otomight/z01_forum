package handlers

import (
	"fmt"
	"forum/internal/config"
	db "forum/internal/database"
	"forum/internal/server/models"
	"forum/internal/utils"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func AddCommentHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
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

	session, ok := r.Context().Value(config.SessionKey).(*db.UserSession)
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

	err = db.AddComment(postID, session.UserID, content)
	if err != nil {
		log.Printf("Failed to add comment: %v", err)
		http.Error(w, "Failed to add comment", http.StatusInternalServerError)
		return
	}
	log.Println("Comment added successfully, redirecting to home page.")
	http.Redirect(w, r, fmt.Sprintf("/post/view/%d", postID), http.StatusSeeOther)
}


func DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	var	session		*db.UserSession
	var	form		models.DeleteCommentForm
	var	comment		*db.Comment
	var	commentID	int
	var	err			error

	if r.Method != http.MethodPost {
		http.Error(w, "Method not alowed", http.StatusMethodNotAllowed)
		return
	}
	if err = utils.ParseForm(r, &form, config.MultipartMaxMemory); err != nil {
		http.Error(
			w, "Unable to parse form:" + err.Error(), http.StatusBadRequest,
		)
		return
	}
	if commentID, err = strconv.Atoi(form.CommentID); err != nil {
		http.Error(
			w, "Failed to delete comment", http.StatusInternalServerError,
		)
		return
	}
	session, _ = r.Context().Value(config.SessionKey).(*db.UserSession)
	comment, _ = db.GetCommentByID(0, commentID)
	if session == nil || comment == nil || session.UserID != comment.AuthorID {
		http.Error(
			w, "You cannot delete this comment!", http.StatusUnauthorized,
		)
		return
	}
	if err = db.DeleteComment(commentID); err != nil {
		http.Error(
			w, "Failed to delete comment", http.StatusInternalServerError,
		)
		return
	}
	if r.Referer() != "" {
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
