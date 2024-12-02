package handlers

import (
	"encoding/json"
	"fmt"
	"forum/internal/config"
	db "forum/internal/database"
	"forum/internal/server/models"
	"net/http"
)

func updateReactionInDb(
	received	models.ReactionPostRequestAjax,
	liked		bool,
) (*models.ReactionPostResponseAjax, error) {
	var	ldl			*db.Reaction
	var	response	models.ReactionPostResponseAjax
	var	err			error

	ldl, err = db.GetReactionByUser(received.PostID, received.UserID)
	if err != nil {
		return nil, err
	}
	if ldl != nil && ldl.Liked == liked {
		err = db.DeleteReaction(received.PostID, received.UserID)
		if err != nil {
			return nil, err
		}
		response.Deleted = true
	} else {
		err = db.AddReaction(received.PostID, received.UserID, liked)
		if err != nil {
			return nil, err
		}
		response.Added = true
		if ldl != nil && ldl.Liked != liked {
			response.Replaced = true
		}
	}
	err = db.UpdatePostReactionsCount(received.PostID)
	return &response, nil
}

func ReactionPostHandler(
	w		http.ResponseWriter,
	r		*http.Request,
	liked	bool,
) {
	var	received	models.ReactionPostRequestAjax
	var	response	*models.ReactionPostResponseAjax
	var	ok			bool
	var	err			error

	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}
	_, ok = r.Context().Value(config.SessionKey).(*db.UserSession)
	if !ok {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	// get data from js
	if err = json.NewDecoder(r.Body).Decode(&received); err != nil {
		http.Error(w, "Invalid request JSON", http.StatusBadRequest)
		return
	}
	// write data in db
	if response, err = updateReactionInDb(received, liked); err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func LikePostHandler(w http.ResponseWriter, r *http.Request) {
	ReactionPostHandler(w, r, true)
}

func DisLikePostHandler(w http.ResponseWriter, r *http.Request) {
	ReactionPostHandler(w, r, false)
}
