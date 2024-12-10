package handlers

import (
	"encoding/json"
	"forum/internal/config"
	db "forum/internal/database"
	"forum/internal/server/models"
	"log"
	"net/http"
)

func updateReactionInDb(
	received	models.ReactionRequestAjax,
	elemType	config.ReactionElemType,
	liked		bool,
) (*models.ReactionResponseAjax, error) {
	var	ldl			*db.Reaction
	var	response	models.ReactionResponseAjax
	var	err			error

	ldl, err = db.GetReactionByUser(elemType, received.ElemID, received.UserID)
	if err != nil {
		return nil, err
	}
	if ldl != nil && ldl.Liked == liked {
		err = db.DeleteUserReaction(
			elemType, received.ElemID, received.UserID,
		)
		if err != nil {
			return nil, err
		}
		response.Deleted = true
	} else {
		err = db.AddReaction(elemType, received.ElemID, received.UserID, liked)
		if err != nil {
			return nil, err
		}
		response.Added = true
		if ldl != nil && ldl.Liked != liked {
			response.Replaced = true
		}
	}
	if err = db.UpdateReactionsCount(elemType, received.ElemID); err != nil {
		return nil, err
	}
	return &response, nil
}

func ReactionHandler(
	w http.ResponseWriter, r *http.Request,
	elemType config.ReactionElemType, liked bool,
) {
	var	received	models.ReactionRequestAjax
	var	response	*models.ReactionResponseAjax
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
	if response, err = updateReactionInDb(received, elemType, liked); err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func PostLikeHandler(w http.ResponseWriter, r *http.Request) {
	ReactionHandler(w, r, config.ReactElemType.Post, true)
}

func PostDisLikeHandler(w http.ResponseWriter, r *http.Request) {
	ReactionHandler(w, r, config.ReactElemType.Post, false)
}

func CommentLikeHandler(w http.ResponseWriter, r *http.Request) {
	ReactionHandler(w, r, config.ReactElemType.Comment, true)
}

func CommentDislikeHandler(w http.ResponseWriter, r *http.Request) {
	ReactionHandler(w, r, config.ReactElemType.Comment, false)
}
