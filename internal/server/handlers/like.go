package handlers

import (
	"encoding/json"
	"fmt"
	"forum/internal/server/models"
	"net/http"
)

func LikeDislikePostHandler(w http.ResponseWriter,
							r *http.Request, liked bool) {
	var		received	models.LikeDislikePostRequestAjax
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
	fmt.Printf("postId=%d userId=%d liked=%v\n",
					received.PostId, received.UserId, liked)
	// err = database.AddLikeDislike(received.PostId, received.UserId, liked)
	// if err  != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	w.WriteHeader(http.StatusOK)
}

func LikePostHandler(w http.ResponseWriter, r *http.Request) {
	LikeDislikePostHandler(w, r, true)
}

func DisLikePostHandler(w http.ResponseWriter, r *http.Request) {
	LikeDislikePostHandler(w, r, false)
}
