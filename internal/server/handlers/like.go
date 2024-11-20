package handlers

import (
	"encoding/json"
	"forum/internal/database"
	"forum/internal/server/models"
	"net/http"
)

func LikeDislikePostHandler(w http.ResponseWriter,
							r *http.Request, liked bool) {
	var		received	models.LikeDislikePostRequestAjax
	var		response	models.LikeDislikePostResponseAjax
	var		err			error
	const	likeCount = 1 // temp define for test
	const	dislikeCount = 1

	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}
	// get data from js
	if err = json.NewDecoder(r.Body).Decode(&received); err != nil {
		http.Error(w, "Invalid request JSON", http.StatusBadRequest)
		return
	}
	// write data in db
	err = database.AddLikeDislike(received.PostId, received.UserId, liked)
	if err  != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// send response with new likes and dislikes count
	response = models.LikeDislikePostResponseAjax{
		LikeCount: likeCount,
		DislikeCount: dislikeCount,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error encoding response JSON",
						http.StatusInternalServerError)
	}
}

func LikePostHandler(w http.ResponseWriter, r *http.Request) {
	LikeDislikePostHandler(w, r, true)
}

func DisLikePostHandler(w http.ResponseWriter, r *http.Request) {
	LikeDislikePostHandler(w, r, false)
}
