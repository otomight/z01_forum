package posthandlers

import (
	"forum/internal/config"
	db "forum/internal/database"
	"forum/internal/server/models"
	"forum/internal/utils"
	"net/http"
	"strconv"
)

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	var	session	*db.UserSession
	var	form	models.DeletePostForm
	var	err		error
	var	post	*db.Post
	var	postID	int

	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}
	if err = utils.ParseForm(r, &form, config.MultipartMaxMemory); err != nil {
		http.Error(w, "Unable to parse form:"+err.Error(),
			http.StatusBadRequest)
		return
	}
	if postID, err = strconv.Atoi(form.PostID); err != nil {
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}
	session, _ = r.Context().Value(config.SessionKey).(*db.UserSession)
	post, _ = db.GetSimplePostByID(0, postID)
	if session == nil || post == nil || session.UserID != post.AuthorID {
		http.Error(w, "You cannot delete this post!", http.StatusUnauthorized)
		return
	}
	if err = db.DeletePost(postID); err != nil {
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
