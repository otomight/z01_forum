package handlers

import (
	"encoding/json"
	"fmt"
	db "forum/internal/database"
	"forum/internal/server/models"
	"forum/internal/server/services"
	"net/http"
)

func updateLikeDislikeInDb(
			received models.LikeDislikePostRequestAjax, liked bool) error {
	var	err	error
	var	ldl	*db.LikeDislike

	ldl, err = db.GetLikeDislikeByUser(received.PostId, received.UserId)
	if err != nil {
		return err
	}
	if ldl != nil && ldl.Liked == liked {
		err = db.DeleteLikeDislike(received.PostId, received.UserId)
		if err != nil {
			return err
		}
	} else {
		err = db.AddLikeDislike(received.PostId, received.UserId, liked)
		if err != nil {
			return err
		}
	}
	err = db.UpdatePostLikesDislikesCount(received.PostId)
	return nil
}

func LikeDislikePostHandler(w http.ResponseWriter,
							r *http.Request, liked bool) {
	var		received	models.LikeDislikePostRequestAjax
	var		err			error
	var		session		*db.UserSession

	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}
	session, _ = services.GetSession(r)
	if !session.IsLoggedIn {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	// get data from js
	if err = json.NewDecoder(r.Body).Decode(&received); err != nil {
		http.Error(w, "Invalid request JSON", http.StatusBadRequest)
		return
	}
	// write data in db
	if err = updateLikeDislikeInDb(received, liked); err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func LikePostHandler(w http.ResponseWriter, r *http.Request) {
	LikeDislikePostHandler(w, r, true)
}

func DisLikePostHandler(w http.ResponseWriter, r *http.Request) {
	LikeDislikePostHandler(w, r, false)
}
